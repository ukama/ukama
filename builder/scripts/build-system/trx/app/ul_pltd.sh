#!/bin/sh

source /tmp/system.conf

DELAY_SYNC=0x10f00009a2514
GAIN1=60
GAIN2=60
usage_cmd()
{
    echo " "    
   
    echo "USAGE: $0 <mode> <TC=test Case number>"
	echo " " 
	
    echo "       Mode:0 = Prepare the Board for the UL PLTD"
	echo "            1 = Start Test"
	echo "            2 = Stop Test"
}

print_uboot_param()
{
	echo " ****** Printin Uboot Parameters ********* " 
	echo $SERVERIP
	echo $BW                
	echo $LOGDISPPORT       
	echo $LOGPARSERIP       
	echo $PWMREGLOW           
	echo $PWMREGHIGH           
	echo $RXFREQ           
	echo $TXFREQ           
	echo $ATTEN1            
	echo $ATTEN2            
	echo $GAIN1             
	echo $GAIN2    
	echo $DUPLEXMODE	
    
}

mode0()
{
echo "*************** Starting Booting Procedure ***************"

cd -
print_uboot_param

cd /usr/bin

echo "2. Configure RF-IF HAB"
        if [ -z $DUPLEXMODE ];then DUPLEXMODE=FDD;fi      
        echo "[$0] : START Override DSP RF Interface parameters in "$DUPLEXMODE" mode"
        if [ -f /usr/bin/updatephyrfparams.sh ]; then     
                /usr/bin/updatephyrfparams.sh ${DUPLEXMODE} WRDDRCFG
                DSPINIT=`oct-linux-memory -w 4  0x83000000 | awk '/Address/{print $4}'`
                echo "[$0] : DONE Override DSP RF Interface parameters, Init Value: $DSPINIT"
        else                                                                          
                echo "[$0] :Script did not correctly retrieve updatephyrfparams.sh by $CFGLOADBY"
        fi                                                                            
        sleep 1   


### Run RF Driver to configure RFIC
echo "3. Configure RF-IC"
oncpu 0 /usr/bin/cn_rfdriver adi <<  EOF
if /usr/bin/ad9362_init_zen${BW}Mhz.txt d ${TXFREQ} u ${RXFREQ} a 1 ${ATTEN1} a 2 ${ATTEN2} g 1 ${GAIN1} g 2 ${GAIN2} gpio 8 q
EOF
sleep 1

sh /etc/init_1pps.sh
sleep 5

## START log #############
echo "5. Start logdisp:-p ${LOGDISPPORT} -a ${LOGPARSERIP}"
oncpu 0 /usr/bin/lsmLogDisp -p ${LOGDISPPORT} -a ${LOGPARSERIP} &
cd /etc

# ms 7ffeff00 6 -w
/usr/bin/oct-linux-memory -w 4 -c 8 0x7ffeff00 6
rm -rf /tmpfs/*
umount -l /tmpfs
mount -t tmpfs -o size=256M /dev/shm /tmpfs
cd /tmpfs

cd /tmpfs

echo "7. Get test vector from TFTP Server ..."
#cp /usr/share/etm/CAL_ETM_TV_${BW}Mhz.tgz .
#cp /usr/share/etm/CAL_UL_TV_${BW}Mhz.tgz .
#tar -zxvf CAL_ETM_TV_${BW}Mhz.tgz
#tar -zxvf CAL_UL_TV_${BW}Mhz.tgz
tftp -g -r TV.tgz ${SERVERIP}
tar -zxvf TV.tgz
sleep 1

## get PLT load
echo "7. Copy PLTD Application to /tmpfs"
cd /usr/bin
gzip -d pltD.gz
cp pltD /tmpfs

## download DSP image
echo "9. Program DSP Image"
cd /tmpfs
oncpu 0 ./pltD -ld normal -stall 2

sleep 2

echo "3. Uncompress Stack obtained from $SWLOADBY"
# already done in bringup
 cd /

 export LD_LIBRARY_PATH=/lib/
 sleep 2
 cd /tmpfs


## display BLER
ccli << EOF
drvCmd -u mask 800
exit
EOF


}


mode1()
{
## execute & stop PLT Test commands
cd /tmpfs

echo "0 /usr/bin/pltD -drv ul -bw $BW -tc ${tc} -c 1"
 oncpu 0 ./pltD -drv ul -bw $BW -tc ${tc} -c 1
cd /usr/bin
tftp -g -r oct-memread ${SERVERIP}
chmod +x oct-memread
export OCTEON_REMOTE_PROTOCOL=linux
 echo "Read the RAW CRC PASS/FAIL Info"
 #/usr/bin/oct-linux-memory -w 4 -c 120 0x10f0000a2b7d0
 echo "CRC PASS oct-memread -w 4 -c 1 0x10f0000a2c234"
 sleep 10
 ./oct-memread  -w 4 -c 1 0x10f0000a2c234
 echo "CRC Fail oct-memread -w 4 -c 1 0x10f0000a2cc64"
 ./oct-memread  -w 4 -c 1 0x10f0000a2cc64
 sleep 5
 
 echo "2. Test running...."
}



mode2()
{
cd /tmpfs

echo "***************** STOP The TEST *********************"
echo "			oncpu 0 ./pltD -drv stop			"
oncpu 0 ./pltD -drv stop

}
total_cmd=$#
file_name=$0
mode=$1
if [ "$total_cmd" -lt 1 ]
then
    usage_cmd
	exit
fi


if [ "$mode" -eq 0 -a "$total_cmd" -ne 1 ] 
then
	usage_cmd
	echo "       Below is the Command to run for mode 0"
	echo "$file_name 0 TV"
	exit
fi

if [ "$mode" -eq 1 -a "$total_cmd" -ne 2 ] 
then
	usage_cmd
	echo "       Command to Start the test case 901"
	echo "$file_name 1 901"
	exit
fi

if [ "$mode" -gt 2 ] 
then
	usage_cmd
	echo ""
	echo "Enter Valid Mode"
	echo ""
	
	exit
fi

case $mode in
		0) TV_NAME=$2; mode0 ;;
		1) tc=$2; mode1 ;;
		2) mode2;;
		*) usage_cmd ; echo "Enter a valid mode !"
esac







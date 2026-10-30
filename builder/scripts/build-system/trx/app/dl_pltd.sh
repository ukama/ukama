#!/bin/sh
source /tmp/system.conf

TV_NAME=ETM_${BW}MHz_TV.tgz
DELAY_SYNC=0x10f00009a2040

usage_cmd()
{
    echo " "    
   
    echo "USAGE: $0 <mode> <TC=test Case number>"
	echo " " 
	
    echo "       Mode:0 = Prepare the Board for the DL PLTD"
	echo "            1 = Start Test"
	echo "            2 = Stop Test"
}

print_system_param()
{
	echo " ****** Printing System Parameters ********* " 
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
	echo $TV_NAME    
}

mode0()
{
echo "*************** Starting Booting Procedure ***************"

print_system_param

echo "1. Configure RF-IF INIT Interface Message in DDR"
#config_rf_if_hab
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
echo "2. Configure RF-IC"
echo "		2.1 if /usr/bin/ad9362_init_zen${BW}Mhz.txt d ${TXFREQ} u ${RXFREQ} a 1 ${ATTEN1} a 2 ${ATTEN2} g 1 ${GAIN1} g 2 ${GAIN2} gpio 8"

oncpu 0 /usr/bin/cn_rfdriver adi <<  EOF
if /usr/bin/ad9362_init_zen${BW}Mhz.txt d ${TXFREQ} u ${RXFREQ} a 1 ${ATTEN1} a 2 ${ATTEN2} g 1 ${GAIN1} g 2 ${GAIN2} gpio 8 q
EOF
sleep 5

echo "3. Configure PTP Module:${MIO_PTP_CLOCK_COMP_VAL}"
sh /etc/init_1pps.sh
sleep 1

## START log #############
## Run /tftpboot/logParser -p9991 in linux host terminal
echo "4. Start logdisp:-p ${LOGDISPPORT} -a ${LOGPARSERIP}"
cd /usr/bin
oncpu 0 lsmLogDisp -p ${LOGDISPPORT} -a ${LOGPARSERIP} &


# ms 7ffeff00 6 -w
/usr/bin/oct-linux-memory -w 4 -c 8 0x7ffeff00 6
rm -rf /tmpfs/*
umount -l /tmpfs
mount -t tmpfs -o size=256M /dev/shm /tmpfs
cd /tmpfs

## get UL PLT test Configuration files
echo "5. Get PLTD Test Vector:${TV_NAME}"
tftp -g -r ${TV_NAME} ${SERVERIP}
tar xzvf ${TV_NAME}

## get PLT load

echo "6. Get PLTD Application from Server:pltD"
 cd /usr/bin
 gzip -d pltD.gz
 cp pltD /tmpfs
 
## get DSP image
echo "7. Get DSP Image"
 cp /etc/dsp.tgz /tmpfs/.
 cd /tmpfs
 tar xzf dsp.tgz
  
sleep 5

## download DSP image
echo "8. Program DSP Image"
 oncpu 0 ./pltD -ld normal -stall 2

 sleep 2

 export LD_LIBRARY_PATH=/lib/

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
#export OCTEON_REMOTE_PROTOCOL=linux
echo "	./pltD -drv dl -bw $BW -tc ${tc} -etm ${etm} -ant ${ant} -c 1   "
 oncpu 0 ./pltD -drv dl -bw $BW -tc ${tc} -etm ${etm} -ant ${ant} -c 1

}



mode2()
{
	cd /tmpfs
#	export OCTEON_REMOTE_PROTOCOL=linux
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
	echo "$file_name 0"
	exit
fi

if [ "$mode" -eq 1 -a "$total_cmd" -ne 4 ] 
then
	usage_cmd
	echo "       Command to Start the test case 901"
	echo "$file_name 1 <testcase number> <etm mode> <antennna mode>"
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
		0) mode0;;
		1) tc=$2; etm=$3;ant=$4;mode1 ;;
		2) mode2;;
		*) usage_cmd ; echo "Enter a valid mode !"
esac







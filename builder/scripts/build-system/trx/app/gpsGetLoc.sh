#!/bin/sh
source /tmp/system.conf

## Will gather 10 seconds worth of seconds by default
if [ -z $GPSDATACOUNT ];then GPSDATACOUNT=10;fi

#GPIOPATH=/sys/class/gpio;

#echo $1 > ${GPIOPATH}/unexport ;
#echo $1 > ${GPIOPATH}/export ;
#echo out > ${GPIOPATH}/"gpio$1"/direction ;
#echo 1 > ${GPIOPATH}/"gpio$1"/active_low ;
checksum_ok=-1
sentce_syntax_ok=-1
pid=-1

config_gps(){

	if [ "${SYNCBY}" == "gps" ];then
	
		/mnt/app/gpsConfigure
		sleep 4;
		echo "[$0] GPS Serial port configured."
		echo 
		
	fi
	
	if [ $pid -ne -1 ];then

		cat /dev/ttyS1 > /tmp/gpslog_raw.txt &
		sleep 3 
					
		pid=`ps  |grep -m 1 /dev/ttyS1 | awk '{print $1}'`
		echo "[$0] Restarted Capture with $pid."
		echo 
	fi	
}

#########################################################################
# traverse gps satellite data			                        #
#########################################################################
literal_gpgsv_str(){
prn_indx=5
elev_indx=6
azim_indx=7
snr_indx=8

while true 
do
	prn=`echo $1 |awk  -v i="$prn_indx" -F "," '{print $i}'`
	elev=`echo $1 |awk  -v i="$elev_indx" -F "," '{print $i}'`
	azim=`echo $1 |awk  -v i="$azim_indx" -F "," '{print $i}'`
	snr=`echo $1 |awk  -v i="$snr_indx" -F "," '{print $i}'`

	if [ -z "$snr" ];then snr="NA"; fi
	if [ -z "$elev" ];then elev="NA"; fi
	if [ -z "$azim" ];then azim="NA"; fi
	
	if [ -n "$prn" ];then
		echo "satellite PRN:$prn  ELEVATION:$elev AZIMUTH:$azim  SNR(dBHz):$snr" >> /tmp/gps_strength;
	else
		break;	
	fi	
	true $((prn_indx+=4))
	true $((snr_indx+=4))
	true $((elev_indx+=4))
	true $((azim_indx+=4))
	
done
}

#########################################################################
# only check sentence number is in order or not				#
#########################################################################
check_syntax(){
	sentce_no=`echo $1 |awk -F "," '{print $3}'`
	if  [ $sentce_no == $2 ]; then 
		sentce_syntax_ok=0;
	else 
		sentce_syntax_ok=-1;
	fi
}
########################################################################
#show satellites in used                                               #
########################################################################
show_satellites_in_used(){
while true
do

	if [ "${SYNCBY}" == "ptp" ];then
		gpgsa_str=`cat /tmp/gpslog_raw.txt |grep "GPS TOD Input:" |awk -F "$" '{print $2}' | grep -a -m 1 "G*GSA"`
		echo "[$0] Grabbing gps location - PTP"
	fi
	if [ "${SYNCBY}" == "gps" ];then
        	gpgsa_str=`cat /tmp/gpslog_raw.txt |grep -a -m 1 "G*GSA"`;
		echo "[$0] Grabbing gps location - gps"
	fi
	checksum_check ${gpgsa_str};
	if [ $checksum_ok == -1 ];then continue;
	else
		break;
	fi
done
	#echo $gpgsa_str >> /tmp/log;
	echo "================================================" >> /tmp/gps_strength
	echo "the satellites used for fix location" >> /tmp/gps_strength
	echo "================================================" >> /tmp/gps_strength
	echo $gpgsa_str | grep G*GSA |awk -F ","        \
                        '{print "PRN: "$4" "$5" "$6" "$7" "$8" "$9" "$10" "$11" "$12" "$13" "$14" "$15}' >> /tmp/gps_strength
}
#########################################################################
# show gps satellites information	                                #
#########################################################################
show_satellites_status(){
while true
do	
	sentce_order=1
	if [ "${SYNCBY}" == "ptp" ];then
                 gpgsv_sentence=`cat /tmp/gpslog_raw.txt |awk -F "$" '{print $2}' | grep -a -m 1 "GPGSV"`
                 sentce_num=`echo $gpgsv_sentence | grep -m 1 "GPGSV" | awk -F "," '{printf $2}'`
		 echo "[$0] Show satellite - ptp"
        fi
        if [ "${SYNCBY}" == "gps" ];then
                gpgsv_sentence=`cat /tmp/gpslog_raw.txt | grep -m 1 "GPGSV"`
                sentce_num=`echo $gpgsv_sentence | grep -m 1 "GPGSV" | awk -F "," '{printf $2}'`
		echo "[$0] Show satellite - gps"
        fi
	sleep 0.1
	if [ "${SYNCBY}" == "ptp" ];then
		cat /tmp/gpslog_raw.txt |awk -F "$" '{print $2}' | grep -a -m $sentce_num "GPGSV" > /tmp/log 
		actual_sentce_num=`wc -l /tmp/log |awk -F " " '{print $1}'`
		#echo actual_sentce_num:$actual_sentce_num;
		if [ $sentce_num != $actual_sentce_num ];then echo "satellites list in view missing data!!"; gather_gps_data; fi
	fi
	if [ "${SYNCBY}" == "gps" ];then
		cat /tmp/gpslog_raw.txt |grep -a -m $sentce_num "GPGSV" > /tmp/log
		actual_sentce_num=`wc -l /tmp/log |awk -F " " '{print $1}'`
                #echo actual_sentce_num:$actual_sentce_num;
		if [ $sentce_num != $actual_sentce_num ];then echo "satellites list in view missing data!!"; gather_gps_data; fi
	fi
	echo "================================================" > /tmp/gps_strength
	satellites_in_view=`cat /tmp/log | grep -m 1 "GPGSV" | awk -F "," '{print $4}'`	
	echo "number of satellites in view:$satellites_in_view" >> /tmp/gps_strength	
	echo "================================================" >> /tmp/gps_strength	
	#checksum check and sentence syntax check for each sentence
	for i in $(cat /tmp/log);
	do
		gpgsv_str=`echo $i |awk -F "GPGSV" '{print $2}'| awk -F "*" '{print $1}'`
		check_syntax $i $sentce_order
		checksum_check $i
		if [ $checksum_ok == -1 ] ||  [ $sentce_syntax_ok == -1 ] ;then echo "bad sentence!!"; gather_gps_data; cat /dev/null > /tmp/gps_strength;break; fi
		literal_gpgsv_str $gpgsv_str
		#echo $gpgsv_str
		true $((sentce_order++))
	done
	show_satellites_in_used;	
	if [ $checksum_ok == 0 ] && [ $sentce_syntax_ok == 0 ]; then cat /tmp/gps_strength;break; fi 
done	
}
#################################################################################
#										#	
# 	NMEA	Decimal								#
#	latitude	0302.78469	03 + (02.78469/60) = 3.046412		#
#	longitude	10141.82531	101 + 41.82531/60) = 101.6971		#
#										#
#################################################################################			
google_api_convert(){
	tmp=`echo $1 |awk -F "." '{print $1}'`
	degree=$(($tmp / 100))
	min=`dc $tmp 100 % p`
	tmp=`echo $1 |awk -F "." '{print $2}'|awk -F "," '{print $1}'`
	point=0.$tmp
	min=`dc $min $point + p`
        #echo $degree, $min $point
        tmp1=`dc $min 60 / p`
        google_format=`dc ${degree} $tmp1 + p`
        dir=`echo $1 |awk -F "," '{print $2}'`
	if [ $dir == 'S' ] || [ $dir == 'W' ];then dir="-"; else dir=""; fi
        echo $dir$google_format
}
#################################################################################
#                                                                               #
#       check sum is deirved by                                                 #
#       bitwise XOR between "$" and "*", including ","                          #
#	0b0000000     0								#
#	0b1000111     G								#
#	----------								#
#	0b1000111     XOR output						#
#	0b1010000     P								#
#	----------								#
#	0b1000111     G								#
#	----------								#
#	0b1010000     XOR output						#
#	0b1000111     G								#
#	----------								#
#	0b0010111     XOR output						#
#	0b1000001     A								#
#	----------								#
#	0b1010110     XOR output						#
#	0b0101100     ,								#
#                                                                               #
#################################################################################
checksum_check(){
	i=1
	if [ "${SYNCBY}" == "ptp" ];then
		nmea_str=`echo $1 |awk -F"*" '{print $1}'`
		echo "[$0] NMEA string - ptp"
	fi
	if [ "${SYNCBY}" == "gps" ];then
        	nmea_str=`echo $1 |awk -F"*" '{print $1}'| awk -F "$" '{print $2}'`
		echo "[$0] NMEA string - gps"
	fi
        ori_sum=`echo $1 |awk -F"*" '{print $2}'`
        #echo ori_sum:$ori_sum                          #debug
        #echo "nmea_str=$nmea_str"                      #debug
        nmea_str_len=${#nmea_str}
        #echo "nmea_length:$nmea_str_len"               #debug
        checksum=`printf %d \"${nmea_str:0:1}\"`

        while true
        do
                more_char=`printf %d \"${nmea_str:$i:1}\"`
                #echo new[${nmea_str:$i:1}]=$more_char  #debug
                checksum=`dc $checksum $more_char eor p`
                #echo checksum=$checksum                #debug
                true $(( i++ ))
                if [ $i -eq $nmea_str_len ];then break;fi
        done
        checksum=`printf %0.2X $checksum`
        #checksum=1                                     #debug
        #echo checksum=$checksum                        #debug
        if [ $ori_sum == $checksum ];then checksum_ok=0;else checksum_ok=-1;fi
}

set_system_time(){
	
	time=`echo $1 |awk -F "," '{print $2}'|awk -F "." '{print $1}'`
	date=`echo $1 |awk -F "," '{print $10}'`
	MM=`echo ${date:2:2}`
	DD=`echo ${date:0:2}`
	YY=`echo 20${date:4:2}`
	
	hh=`echo ${time:0:2}`
	mm=`echo ${time:2:2}`
	ss=`echo ${time:4:2}`
	echo "[$0] set system time:$YY-$MM-$DD $hh:$mm:$ss"
	date -s "$YY-$MM-$DD $hh:$mm:$ss"
	
}

gps_exit(){
    if [ "${SYNCBY}" == "ptp" ];then
       echo -e "\nptp2 e d off\n" > /dev/ttyS1                                      
    fi
    sleep 1
    exit 1
}

gather_gps_data(){

#config_gps

missing_counter=0

while [ $missing_counter -lt 10 ];
do
       echo -n "[$0] Gather GPS DATA"
       
       if [ "${SYNCBY}" == "ptp" ];then
		cat /dev/ttyS1 |grep "GPS TOD Input:" > /tmp/gpslog_raw.txt &
       elif [ "${SYNCBY}" == "gps" ];then
		cat /dev/ttyS1 > /tmp/gpslog_raw.txt &
       else
		break	
       fi			

       pid=`ps  | grep -m 1 /dev/ttyS1 |awk '{print $1}'`
       
       i=0
       echo  
       echo "[$0] Started NMEA capture at $pid." 
       while  true ;do                                                                              
          if [ $i -eq $GPSDATACOUNT ];then break;fi                                                 
             let i=$i+1                                                                             
             echo -n " ."                                                                           
             sleep 1                                                                                
          done                                                                                          
          
       kill -9 $pid	       
       echo
               
       gpgsv_present=`cat /tmp/gpslog_raw.txt |grep -i -c gpgsv`                                    
       gprmc_present=`cat /tmp/gpslog_raw.txt| grep -i -c g*rmc`
       gpgsa_present=`cat /tmp/gpslog_raw.txt| grep -i -c g*gsa`
       
       if [ $gpgsv_present -eq 0 ];then 
           echo "[$0] Incomplete gps data GPGSV"
           let missing_counter=$missing_counter+1
           #config_gps
           continue
       elif [ $gprmc_present -eq 0 ];then 
          echo "[$0] Incomplete GPS data G*RMC"
          let missing_counter=$missing_counter+1
          #config_gps
          continue
       elif [ $gpgsa_present -eq 0 ];then 
           echo "[$0] Incomplete GPS data GPGSA"
           let missing_counter=$missing_counter+1
           #config_gps
           continue
       else
       	   break
       fi
    
done

if [ $missing_counter -eq 10 ];then echo "GPS data gather failed - Exiting"; gps_exit;fi     

}

trap ctrl_c INT

ctrl_c(){
    echo "[$0] Trapped CTRL-C"
    gps_exit
}
    
#################################################################################################
#                                                                                               #
#                                                                                               #
#       RMC - NMEA has its own version of essential gps pvt (position, velocity, time) data.    #
#       It is called RMC, The Recommended Minimum, which will look similar to:                  #
#                                                                                               #
#                                                                                               #
#               123519       Fix taken at 12:35:19 UTC                                          #
#               A            Status A=active or V=Void.                                         #
#               4807.038,N   Latitude 48 deg 07.038' N                                          #
#               01131.000,E  Longitude 11 deg 31.000' E                                         #
#               022.4        Speed over the ground in knots                                     #
#               084.4        Track angle in degrees True                                        #
#               230394       Date - 23rd of March 1994                                          #
#               003.1,W      Magnetic Variation                                                 #
#               *6A          The checksum data, always begins with *                            #
#                                                                                               #
#################################################################################################
#show_satellites_status 

####
#Check if GPS is enabled                                                 
##
if [ -z ${SYNCBY} ]; then SYNCBY=internal; echo "[$0] SYNCBY : undeclared setting to internal";fi
if [ -z $GPSLOCKDEBUG ];then GPSLOCKDEBUG=0;fi

# Make sure logs are off first
if [ "${SYNCBY}" == "ptp" ];then
        echo -e "\nptp2 e d off\n" > /dev/ttyS1              
        sleep 1
        #All we need is TOD signal status to be stable
        GPSPTPENABLE=`/usr/bin/ptprd.sh 'ptp2 g s' 'TOD signal' |grep -c ' stable'`
else
	GPSPTPENABLE=0
fi
  
  
if [ $GPSLOCKDEBUG -gt 0 ];then
	echo "SYNCBY = $SYNCBY"
	echo "GPSPTPENABLE = $GPSPTPENABLE"                              
fi
                                                                                                                                                   
if [ ! "${SYNCBY}" == "gps" -a ! $GPSPTPENABLE -gt 0 ];then echo "GPS not present or locked.  EXIT"; gps_exit;fi
                                                                                                                                                                                  
#Setup PTP module to check GPS log on TOD inout
if [ "$GPSPTPENABLE" == "1" ];then
    if [ "$GPSLOCKDEBUG" == "1" ];then
	echo "Sending ptp2 debug commands..."
    fi
    echo -e "\nptp2 e d on\n\nptp2 e v 19\n" > /dev/ttyS1
fi
####
re=0
while true
do
	Latitude=""
	Longitude=""
	#echo 1 > ${GPIOPATH}/"gpio$1"/value;

	#show_satellites_status 
	#check gps fix or not
	## When SYNCBY is set to ptp only
        gather_gps_data

	if [ "${SYNCBY}" == "ptp" ];then
		gps_fix_str=`cat /tmp/gpslog_raw.txt |awk -F "$" '{print $2}' | grep -a -m 1 "G*GSA"`
	fi
	if [ "${SYNCBY}" == "gps" ];then
        	gps_fix_str=`cat /tmp/gpslog_raw.txt |grep -a -m 1 "G*GSA"`;
	fi
        checksum_check ${gps_fix_str};
        
        if [ $checksum_ok == -1 ];then continue;        \
                                                        \
        else
                GPSFIX=`echo $gps_fix_str|awk -F"," '{print $3}'`;      \
                echo "gps location is fixed=$GPSFIX";
        fi
        
        #capture data, time, date, longitude, lantitude
	if [ "${SYNCBY}" == "ptp" ];then
		rmc_str=`cat /tmp/gpslog_raw.txt |awk -F "$" '{print $2}' | grep -a -m 1 "G*RMC"`
	fi

        if [ "${SYNCBY}" == "gps" ];then
		rmc_str=`cat /tmp/gpslog_raw.txt |grep -a -m 1 "G*RMC"`
	fi
        checksum_check ${rmc_str}
        if [ $checksum_ok == -1 ];then continue;        \
							\
        else
                Latitude=`echo $rmc_str |awk -F"," '{print $4}'`,`echo $rmc_str |awk -F"," '{print $5}'`;
                Longitude=`echo $rmc_str |awk -F"," '{print $6}'`,`echo $rmc_str|awk -F"," '{print $7}'`;
                
        fi
        echo "GPGSA:$gps_fix_str" > /tmp/gps_log
        echo "G*RMC:$rmc_str" >> /tmp/gps_log
        echo "latitude:$Latitude, logngitude:$Longitude" >> /tmp/gps_log
	
	if [ "$GPSSHOWSTATUS" == "1" ];then	
		show_satellites_status;
	fi	
        if [ $GPSFIX != "1" ] && [ -n ${Latitude} ] && [ -n $Longitude ];then   \
                google_lat=`google_api_convert  ${Latitude}`;
                google_long=`google_api_convert  ${Longitude}`;
                echo ${google_lat},${google_long} > /etc/gps_loc.log;
                echo 
                echo "[$0] GPS location :${google_lat},${google_long}";
                ## If SYNCBY is set to ptp then time is already updated by PTP module.
                if [ "${SYNCBY}" == "gps" ];then
		    set_system_time ${rmc_str};
                fi
                #echo 0 > ${GPIOPATH}/"gpio$1"/value;
                gps_exit
	else	
		echo -n ". "; 
	fi	
        true $(( re++ ))
	if [ $re -eq 10 ] && [ $GPSFIX != "3" ];then echo "GPS fix location time out!!"; break;fi
        #echo 0 > ${GPIOPATH}/"gpio$1"/value;
        sleep 0.1
done

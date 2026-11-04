#!/bin/sh

source /tmp/system.conf


#-----------------------------------------------------------------------------
# tendency will look at the direction of the change and correct by a fix amount
# absolute will look at the distance and correct error to average
# initial will use absolute initially and then switch to Tendency algo.
algo="initial"  

#------------------------------------------------------------------------------
test_time=5		        # the sleep/wait period
#------------------------------------------------------------------------------
TARGET_REG="BBP_RFIF_SAMPLE_ADJ_ERROR"	# target PHY SAMPLE COUNT register
#------------------------------------------------------------------------------

if [ "$BW" == "5" ];then
    clockspeed=15360		# target #clocks per 10ms period
    clockspeed_div2=7680	# target #clocks per 10ms period DIVIDED by 2
    tolerance=7677		# This is clockspeed/2 +/- 2 effectively.
    pwm_sample_to_unit=384	# N PWM unit = 1 sample
    pwm_adj_step=48		# N step when in tendancy mode
    dac_sample_to_unit=128	# N PWM unit = 1 sample
    dac_adj_step=32		# N step when in tendancy mode
elif [ "$BW" == "10" ];then
    clockspeed=30720		# target #clocks per 10ms period
    clockspeed_div2=15360	# target #clocks per 10ms period DIVIDED by 2
    tolerance=15357		# This is clockspeed/2 +/- 2 effectively.
    pwm_sample_to_unit=256	# N PWM unit = 1 sample
    pwm_adj_step=32		# N step when in tendancy mode
    dac_sample_to_unit=64	# N PWM unit = 1 sample
    dac_adj_step=16		# N step when in tendancy mode
elif [ "$BW" == "15" ];then
    clockspeed=46080		# target #clocks per 10ms period
    clockspeed_div2=23040	# target #clocks per 10ms period DIVIDED by 2
    tolerance=23037		# This is clockspeed/2 +/- 2 effectively.
    pwm_sample_to_unit=128	# N PWM unit = 1 sample
    pwm_adj_step=16		# N step when in tendancy mode
    dac_sample_to_unit=96	# N PWM unit = 1 sample
    dac_adj_step=16		# N step when in tendancy mode
elif [ "$BW" == "20" ];then
    clockspeed=61440		# target #clocks per 10ms period
    clockspeed_div2=30720	# target #clocks per 10ms period DIVIDED by 2
    tolerance=30717		# This is clockspeed/2 +/- 2 effectively.
    pwm_sample_to_unit=64	# N PWM unit = 1 sample
    pwm_adj_step=16		# N step when in tendancy mode
    dac_sample_to_unit=64	# N PWM unit = 1 sample
    dac_adj_step=16		# N step when in tendancy mode
fi

#------------------------------------------------------------------------------
if [ "$TCXO" == "pwm" ];then
sample_to_unit=$pwm_sample_to_unit
adj_step=$pwm_adj_step

else
echo "Selecting DAC control steps"
sample_to_unit=$dac_sample_to_unit
adj_step=$dac_adj_step

fi
#------------------------------------------------------------------------------
##  offset min max moved to conditional statement belo
#------------------------------------------------------------------------------

## For PWM
if [ "$TCXO" == "pwm" ];then
    offset_max=60000  #   ## 16 = 8Hz @2.1Ghz
    offset_min=15000  #

    default_offset=`/usr/bin/printpwm |grep LOW_TIME |awk '{print $3}'`
## For DAC
else 
    offset_max=65535
    offset_min=0
    default_offset=`fprintenv dac_offset -r`	# the offset value in uboot
fi


## Concept of target ofset is erronous as we are trying to avoid calibration
target_offset=${default_offset}

clockvalue=0

noadj=0

#------------------------------------------------------------------------------
DoPingExit=0
pingIP="10.102.81.2"
#------------------------------------------------------------------------------
average_num=10
#------------------------------------------------------------------------------
#
# sub function
#
#------------------------------------------------------------------------------
#
#
set_OFFSET()	# to set the DAC offset
{
	#----------------------------------------------------------------------
	if [ ! -n "$1" ] ; then
		offset=${default_offset}
	else
		offset=$1
	fi
	#----------------------------------------------------------------------
	#if [ "${offset}" -eq "${current_offset}" ] ; then
	#	offset=`expr ${offset} - 1`
	#fi
	#----------------------------------------------------------------------

	


	## For PWM
	if [ "$TCXO" == "pwm" ];then
	    echo -n "Adjusting PWM LOW to: $offset"
	    #oncpu 0 /usr/bin/oct_clk_ctrl -e ${offset} > /dev/null 2>&1
    	    oncpu 0 /usr/bin/cn_rfdriver adi << EOF >> /tmp/xoserverwrite.log 2>&1
    	    pwm l ${offset} r q
EOF

    	else
	##for DAC
	    #oncpu 0 /usr/bin/ad5621Tester ${offset} > /dev/null 2>&1
	    oncpu 0 /usr/bin/cn_rfdriver adi << EOF > /dev/null 2>&1           
	    dax ${offset} q                                                  
EOF
	fi	

	#----------------------------------------------------------------------
	echo $offset > /tmp/last_offset.txt
	#----------------------------------------------------------------------

}
#------------------------------------------------------------------------------
#
#

#------------------------------------------------------------------------------
#------------------------------------------------------------------------------
#
#
get_RFIF_count()	# to get 10ms sample count value
{
	sleep_time=$1
	#----------------------------------------------------------------------
	oct-linux-csr ${TARGET_REG} > /tmp/gps1
	#----------------------------------------------------------------------
	gps1=`cat /tmp/gps1 | grep "OFFSET" | awk '{print $5}'`


	#----------------------------------------------------------------------
	clockvalue=`expr ${gps1} - \( ${clockspeed} / 2 \)`
	#----------------------------------------------------------------------

	sleep $1
}


#------------------------------------------------------------------------------
do_adjust()	# do the DAC||PWM offset adjustment
{
	#----------------------------------------------------------------------
	current_offset=$1
	## $2 is the direction of the change
	#----------------------------------------------------------------------
	echo "[do_adjust] Adjusting : $2 "
	if [ "$2" -lt "0" ];then
	    ## careful offset_min and $2 must be consistent
	    offset=`expr ${current_offset} + $2`
	    echo  "[do_adjust]  $offset > $offset_min ?"
	    if [ "${offset}" -gt "${offset_min}" ];then 
		set_OFFSET ${offset}
	    else
		echo "[do_adjust]  ${current_offset} min offset saturated to ${offset} ==> Unlikely correct"
	    fi
	else
	    offset=`expr ${current_offset} + $2`
	    echo  "[do_adjust]  $offset < $offset_max ?"
	    if [ "${offset}" -lt "${offset_max}" ];then
	        set_OFFSET ${offset}
	    else
	        echo "[do_adjust]  $2 max offset saturated to ${current_offset} ==> Unlikely correct"
	    fi
	fi
}

check_average_clock_value()
{
	#----------------------------------------------------------------------
	num=0
	sum=0
	sign=0
	#----------------------------------------------------------------------
	for i in `cat /tmp/clockvalues | tail -n ${average_num}`; do
		num=`expr ${num} + 1`
		sum=`expr ${sum} + ${i#-}`  #take absolute value
		if [ "${i}" -lt "0" ];then
		    sign=`expr ${sign} - 1`
		else
		    sign=`expr ${sign} + 1`
		fi
	done

	if [ "$TCXO" != "pwm" ];then
		sign=`expr -1 \* ${sign}`
	fi

	#----------------------------------------------------------------------
	average=`expr  ${sum} / ${num}`
	#----------------------------------------------------------------------
	current_offset=`cat /tmp/last_offset.txt`
	#----------------------------------------------------------------------
	echo -n "num:[${num}] noadj: [${noadj}] average: [${average}] current offset: ${current_offset}"
	#----------------------------------------------------------------------	
	
	# check how far average is from @10Mhz is 15360 (perfect value when clockspeed = 30720 (10MhzBW)
	distance=`expr ${clockspeed_div2} - ${average}`
	
	if [ "${algo}" == "tendency" ];then
		correction=$adj_step
	elif [ "${algo}" == "absolute" -o "${algo}" == "initial" ];then
		correction=`expr ${distance} \* ${sample_to_unit}`     # @10Mhz, 256 units = 65Hz when pwmreg high ~=32768. 
	fi
        
        echo -n " corr = $correction"                                                              
         
	
	if [ "${num}" -eq "${average_num}" ]; then
		echo -n "+ "
		
		if [ "${average}" -gt "${tolerance}" ] ; then
	            noadj=`expr ${noadj} + 1`
	            ## If we were using Initial convergence algo, switch to Tendency algo now what we are witin limits.
	            #
	            if [ "${algo}" == "initial" ];then echo -e "\n\nWithin Tolerence, switching to Tendency Algo\n"; algo=tendency; fi
	            #
		elif [ "${sign}" -lt "0" ] ; then
		    do_adjust ${current_offset} ${correction}
		    noadj=0
		elif [ "${sign}" -gt "0" ] ; then
			do_adjust ${current_offset} -${correction}
			noadj=0
		fi
		
		## After an absolute or initial adjustement, restart averaging filter.
		if [ "${algo}" != "tendency" -a "${noadj}" == "0" ];then rm -rf /tmp/clockvalues; sleep 5;fi
	fi
	
	
	#----------------------------------------------------------------------

	#----------------------------------------------------------------------
	echo ""
	#----------------------------------------------------------------------
}

check_environment()
{

    #-------------------------------------------------------------------
    # check the last_offset
    echo "Create /tmp/last_offset.txt"
    echo "${default_offset}" > /tmp/last_offset.txt
    #----------------------------------------------------------------------
}

#------------------------------------------------------------------------------
#
# main function
#
#------------------------------------------------------------------------------
check_environment
#------------------------------------------------------------------------------
rm -rf /tmp/clockvalues
#------------------------------------------------------------------------------

#------------------------------------------------------------------------------
echo "start to check/observer clock value"
#------------------------------------------------------------------------------
loop=1
while [ "1" ] ; do
	#----------------------------------------------------------------------
	get_RFIF_count ${test_time}
	#----------------------------------------------------------------------
	echo ${clockvalue} >> /tmp/clockvalues
	#----------------------------------------------------------------------
	echo -n "[`date +%Y/%m/%d_%H:%M:%S`] [${loop}] clock: [${clockvalue}] "
	echo -n "[`date +%Y/%m/%d_%H:%M:%S`] [${loop}] clock: [${clockvalue}] " >> /tmp/xoserverwrite.log	
	#----------------------------------------------------------------------
	check_average_clock_value
	#----------------------------------------------------------------------

	#----------------------------------------------------------------------
	loop=`expr ${loop} + 1`
	#----------------------------------------------------------------------
done
#------------------------------------------------------------------------------

#!/bin/bash
#Copyright (c) 2021-present, Ukama.

SYSDIR=/tmp/sys
SYSFSDIRHWMON=/tmp/sys/class/hwmon/hwmon0/
SYSFSDIRGPIO=/tmp/sys/class/gpio
SYSFSLED=/tmp/sys/class/led/
RANDOM=$$
I2CDEVPATH=/tmp/sys/bus/i2c/devices

create_sysfs_for_eeprom() {
	I2CBUS=$1
	I2CADDRESS=$2
	UNITINFO=$3
	ACTINGMASTER=$4
	EPATH=${I2CDEVPATH}/i2c-${I2CBUS}/${I2CBUS}-${I2CADDRESS}/
	mkdir -p ${EPATH};
	cd ${EPATH}
	#touch eeprom
	echo "Acting master is ${ACTINGMASTER}"
	if [ "${ACTINGMASTER}" -eq 1 ] ; then
		ln -s ${EPATH}eeprom /tmp/sys/${UNITINFO}-systemdb
	fi
}

create_sysfs_for_tmp() {
	COUNT=$1
	FTEMPVALUE=temp${COUNT}_input
	FMINVALUE=temp${COUNT}_min
	FMAXVALUE=temp${COUNT}_max
	FCRITVALUE=temp${COUNT}_crit
	FCRITHYST=temp${COUNT}_crit_hyst
	FMAXHYST=temp${COUNT}_max_hyst
	FMINALARM=temp${COUNT}_min_alarm
	FMAXALARM=temp${COUNT}_max_alarm
	FCRITALARM=temp${COUNT}_crit_alarm
	FOFFSET=temp${COUNT}_offset

	touch ${FTEMPVALUE}
	[ -f "${FTEMPVALUE}" ] && { echo "${FTEMPVALUE} created.";}
	echo 45000 > ${FTEMPVALUE}
	echo "Reading ${FTEMPVALUE}" `cat ${FTEMPVALUE}`

	touch ${FMINVALUE}
	[ -f "${FMINVALUE}" ] && { echo "${FMINVALUE} created.";}
	echo -25000 > ${FMINVALUE}
	echo "Reading ${FMINVALUE}" `cat ${FMINVALUE}`

	touch ${FMAXVALUE}
	[ -f "${FMAXVALUE}" ] && { echo "${FMAXVALUE} created.";}
	echo 75000 > ${FMAXVALUE}
	echo "Reading ${FMAXVALUE}" `cat ${FMAXVALUE}`

	touch ${FCRITVALUE}
	[ -f "${FCRITVALUE}" ] && { echo "${FCRITVALUE} created."; }
	echo 85000 > ${FCRITVALUE}
	echo "Reading ${FCRITVALUE}" `cat ${FCRITVALUE}`

	touch ${FCRITHYST}
	[ -f "${FCRITHYST}" ] && { echo "${FCRITHYST} created."; }
	echo 2000 > ${FCRITHYST}
	echo "Reading ${FCRITHYST}" `cat ${FCRITHYST}`

	touch ${FMAXHYST}
	[ -f "${FMAXHYST}" ] && { echo "${FMAXHYST} created."; }
	echo 2000 > ${FMAXHYST}
	echo "Reading ${FMAXHYST}" `cat ${FMAXHYST}`

	touch ${FOFFSET}
	[ -f "${FOFFSET}" ] && { echo "${FOFFSET} created."; }
	echo 5000 > ${FOFFSET}
	echo "Reading ${FOFFSET}" `cat ${FOFFSET}`

	touch ${FMINALARM}
	[ -f "${FMINALARM}" ] && { echo "${FMINALARM} created."; }
	echo 0 > ${FMINALARM}
	echo "Reading ${FMINALARM}" `cat ${FMINALARM}`

	touch ${FMAXALARM}
        [ -f "${FMAXALARM}" ] && { echo "${FMAXALARM} created."; }
        echo 0 > ${FMAXALARM}
        echo "Reading ${FMAXALARM}" `cat ${FMAXALARM}`

	touch ${FCRITALARM}
        [ -f "${FCRITALARM}" ] && { echo "${FCRITALARM} created."; }
        echo 0 > ${FCRITALARM}
        echo "Reading ${FCRITALARM}" `cat ${FCRITALARM}`

}

create_sysfs_for_adt7481() {
	SYSDIRPATH=$1
	CURDEVNO=$2
	cd ${SYSDIRPATH};
	mkdir -p adt7481_${CURDEVNO};
	cd adt7481_${CURDEVNO};
	CURPWD=`pwd`
	echo "Currently working in ${CURPWD}."
	echo "Creating Sysfs For ADT7481";
	for (( iter=1; iter<=3 ; iter++ ))
	do
		create_sysfs_for_tmp ${iter}
	done
}

create_sysfs_for_tmp464() {
        SYSDIRPATH=$1
	CURDEVNO=$2
        cd ${SYSDIRPATH};
	mkdir -p tmp464_${CURDEVNO};
        cd tmp464_${CURDEVNO};
        CURPWD=`pwd`
        echo "Currently working in ${CURPWD}."
        echo "Creating Sysfs For TMP464";
        for (( iter=1; iter<=3 ; iter++ ))
        do
                create_sysfs_for_tmp ${iter}
        done
}

create_sysfs_for_se98() {
	SYSDIRPATH=$1
	CURDEVNO=$2
	cd ${SYSDIRPATH};
	mkdir -p se98_${CURDEVNO}
	cd se98_${CURDEVNO}
	CURPWD=`pwd`
        echo "Currently working in ${CURPWD}."
        echo "Creating Sysfs For SE98";
        for (( iter=1; iter<=1 ; iter++ ))
        do
                create_sysfs_for_tmp ${iter}
        done
}

create_sysfs_for_att() {
	SYSDIRPATH=$1
	NUM=$2
	cd ${SYSDIRPATH};
        mkdir -p att${NUM}
	cd att${NUM}

	touch in0_attvalue
	echo 63 > in0_attvalue

	touch in0_latch
	echo 0 > in0_latch
}

create_sysfs_led() {
	SDIR=$1
	mkdir -p ${SDIR};
        cd ${SDIR};
	touch brightness;
	echo 0 > brightenss;
	touch max_brightness;
	echo 255 > max_brightness;
	touch trigger
	echo "none" > trigger	
}

create_sysfs_for_leds() {
	SYSDIRPATH=${SYSFSLED}
        LEDNUM=$2
	LEDDIR=led${LEDNUM} 
	cd ${SYSDIRPATH};
        mkdir -p ${LEDDIR}
	cd ${LEDDIR}
	CPWD=`pwd`
	cd ${CPWD}
	create_sysfs_led red
	cd ${CPWD}
        create_sysfs_led green
	cd ${CPWD}
	create_sysfs_led blue	
}	


create_sysfs_for_ads1015() {
	SYSDIRPATH=$1
	NUM=$2
	SDIR=adc${NUM}
	cd ${SYSDIRPATH}
	mkdir -p ${SDIR}
	cd ${SDIR}
	touch in0_input
        echo $RANDOM > in0_input
	
	touch in1_input
        echo $RANDOM > in1_input

	touch in2_input
        echo $RANDOM > in2_input

	touch in3_input
        echo $RANDOM > in3_input

	touch in4_input
        echo $RANDOM > in4_input

	touch in5_input
        echo $RANDOM > in5_input

	touch in6_input
        echo $RANDOM > in6_input

	touch in7_input
        echo $RANDOM > in7_input		
}

create_sysfs_for_inpgpio() {
	SYSDIRPATH=${SYSFSDIRGPIO}
	GPIONUM=$1
	cd ${SYSDIRPATH};
	mkdir -p gpio${GPIONUM}
	cd gpio${GPIONUM}
	touch direction 
	echo "in" > direction;
	touch value
	echo 1 > value
	touch edge
	echo "rising" > edge
	touch active_low
	echo 0 > polarity
}

create_sysfs_for_outgpio() {
        SYSDIRPATH=${SYSFSDIRGPIO}
        GPIONUM=$1
        cd ${SYSDIRPATH};
        mkdir -p gpio${GPIONUM}
        cd gpio${GPIONUM}
        touch direction
        echo "out" > direction;
        touch value
        echo 1 > value
        touch edge
        echo "both" > edge
        touch active_low
        echo 0 > polarity
}


create_sysfs_for_ina226() {
	SYSDIRPATH=$1
        CURDEVNO=$2
        cd ${SYSDIRPATH};
        mkdir -p ina226_${CURDEVNO}
        cd ina226_${CURDEVNO}
        CURPWD=`pwd`
        echo "Currently working in ${CURPWD}."
        echo "Creating Sysfs For INA226";
	
	SHUNTVOLTAGE=in0_input
	BUSVOLTAGE=in1_input
	CURRENT=curr1_input
	POWER=power1_input
	SHUNTRESISTOR=shunt_resistor
	CRITLOWSHUNTVOLTAGE=in0_lcrit
	CRITHIGHSHUNTVOLTAGE=in0_crit
	SHUNTVOLTAGECRITLOWALARM=in0_lcrit_alarm
	SHUNTVOLTAGECRITHIGHALARM=in0_crit_alarm
	CRITLOWBUSVOLTAGE=in1_lcrit
	CRITHIGHBUSVOLTAGE=in1_crit
	BUSVOLTAGECRITLOWALARM=in1_lcrit_alarm
	BUSVOLTAGECRITHIGHALARM=in1_crit_alarm
	CRITHIGHPWR=power1_crit
	CRITHIGHPWRALARM=power1_crit_alarm
	UPDATEINTERVAL=update_interval
	
	INASYSFS=(${SHUNTVOLTAGE} ${BUSVOLTAGE} ${CURRENT} ${POWER} ${SHUNTRESISTOR} ${CRITLOWSHUNTVOLTAGE} ${CRITHIGHSHUNTVOLTAGE}  ${CRITLOWBUSVOLTAGE} ${CRITHIGHBUSVOLTAGE} ${CRITHIGHPWR} ${UPDATEINTERVAL} )
	INASYSFS_MIN=("1850" "11700" "4700" "60000" "9500" "1850" "2100" "11700" "12250" "60500" "1500")
        INASYSFS_MAX=("2150" "12300" "5300" "63000" "10500" "1900" "2150" "11750" "12300" "61000" "3000") 	
	INASYSFSVAL=(
		INASYSFS[@]
		INASYSFS_MIN[@]
		INASYSFS_MAX[@]
		)
	CNT=0
	for FILE in "${INASYSFS[@]}"
	do
   		echo ${FILE}
		touch ${FILE}
        	[ -f "${FILE}" ] && { echo "${FILE} created."; }
		echo "File: ${!INASYSFSVAL[0]:${CNT}:1} Value at [${CNT}]: { [${!INASYSFSVAL[1]:${CNT}:1}], [${!INASYSFSVAL[2]:${CNT}:1}] }"
        	VAL=`shuf -i ${!INASYSFSVAL[1]:${CNT}:1}-${!INASYSFSVAL[2]:${CNT}:1} -n 1`
		echo "Writing ${VAL} to ${FILE}"
		echo ${VAL} > ${FILE}
		echo "Reading ${FILE}" `cat ${FILE}`
		CNT=$((CNT+1))	
	done

	INASYSFS_ALARM=(${SHUNTVOLTAGECRITLOWALARM} ${SHUNTVOLTAGECRITHIGHALARM} ${BUSVOLTAGECRITLOWALARM} ${BUSVOLTAGECRITHIGHALARM} ${CRITHIGHPWRALARM} )
	 for FILE in "${INASYSFS_ALARM[@]}"
        do
                echo ${FILE}
                touch ${FILE}
                [ -f "${FILE}" ] && { echo "${FILE} created."; }
                echo 0 > ${FILE}
                echo "Reading ${FILE}" `cat ${FILE}`
        done

}

create_sysfs_for_module() {
	UNIT=$1
	MODID=$2
	MASTER=$3
	cd ${SYSFSDIRHWMON}
	mkdir -p ${MODID}
	cd ${MODID}
	MODDIR=`pwd`
	echo "Creating Sysfs for module ${MODID}"
	case "${MODID}" in
		1)
			create_sysfs_for_eeprom 0 0050 ${UNIT} ${MASTER}
			cd ${MODDIR}	
			create_sysfs_for_tmp464 ${MODDIR} 1 
			cd ${MODDIR}
			create_sysfs_for_tmp464 ${MODDIR} 2
			cd ${MODDIR}
			create_sysfs_for_se98 ${MODDIR} 1
			cd ${MODDIR}
			create_sysfs_for_ina226 ${MODDIR} 1
			cd ${MODDIR}
			create_sysfs_for_ina226 ${MODDIR} 2
			cd ${MODDIR}
			create_sysfs_for_ina226 ${MODDIR} 3
			cd ${MODDIR}
                        create_sysfs_for_leds ${MODDIR} 0
                        cd ${MODDIR}
                        create_sysfs_for_leds ${MODDIR} 1
                        cd ${MODDIR}
                        create_sysfs_for_leds ${MODDIR} 2
                        cd ${MODDIR}
                        create_sysfs_for_leds ${MODDIR} 3
                        cd ${MODDIR}
                        ;;

		2)
			create_sysfs_for_eeprom 1 0050 ${UNIT} ${MASTER}
                        cd ${MODDIR}
			create_sysfs_for_tmp464 ${MODDIR} 1
			cd ${MODDIR}
			create_sysfs_for_tmp464 ${MODDIR} 2
			cd ${MODDIR}
			create_sysfs_for_se98 ${MODDIR} 1
			cd ${MODDIR}
			create_sysfs_for_ina226 ${MODDIR} 1
			cd ${MODDIR}
			create_sysfs_for_ina226 ${MODDIR} 2
			cd ${MODDIR}
			;;
		3)
			create_sysfs_for_eeprom 1 0051 ${UNIT} ${MASTER}
                        cd ${MODDIR}
			create_sysfs_for_adt7481 ${MODDIR} 1
			cd ${MODDIR}
			create_sysfs_for_ina226 ${MODDIR} 1
                	cd ${MODDIR}
			;;
		4)
			create_sysfs_for_eeprom 1 0052 ${UNIT} ${MASTER} 
			cd ${MODDIR}
                        create_sysfs_for_se98 ${MODDIR} 1
                        cd ${MODDIR}
			create_sysfs_for_tmp464 ${MODDIR} 1
                        cd ${MODDIR}
			create_sysfs_for_ads1015 ${MODDIR} 1
                        cd ${MODDIR}
			create_sysfs_for_att ${MODDIR} 1
                        cd ${MODDIR}
			create_sysfs_for_att ${MODDIR} 2
                        cd ${MODDIR}
			create_sysfs_for_inpgpio 38
			cd ${MODDIR}
                        create_sysfs_for_inpgpio 35
			cd ${MODDIR}
                        create_sysfs_for_inpgpio 34
			cd ${MODDIR}
                        create_sysfs_for_outgpio 63
                        cd ${MODDIR}
                        create_sysfs_for_outgpio 61
                        cd ${MODDIR}
                        create_sysfs_for_outgpio 40
			;;
		5)
			create_sysfs_for_eeprom 0 0051 ${UNIT} ${MASTER}
                        cd ${MODDIR}
			create_sysfs_for_se98 ${MODDIR} 1
			cd ${MODDIR}
                        create_sysfs_for_tmp464 ${MODDIR} 1
			cd ${MODDIR}
			create_sysfs_for_leds ${MODDIR} 0
			cd ${MODDIR}
                        create_sysfs_for_leds ${MODDIR} 1
			cd ${MODDIR}
                        create_sysfs_for_leds ${MODDIR} 2
			cd ${MODDIR}
                        create_sysfs_for_leds ${MODDIR} 3
			cd ${MODDIR}
			;;
		*)
			echo "Unknown module number."
			;;
		esac
}

usage() {
	echo "./prepare_env.sh [option]"
	echo "Options:"
	echo " -c | --clean		Clean the sysfs dir."
	echo " -u | --unittype          Create sysfs for unit type.
					Valid unit types: cnode-lte/anode"
}
     	
clean_sysfs_dir() {
	echo "Clean ${SYSDIR}"
	rm -rf ${SYSDIR}
	sync
}

create_sysfs_for_unit() {
CPWD=$1
UNITTYPE=$2
echo "Current Working DIR is ${CPWD}"
case ${UNITTYPE} in
	"cnode-lte")
		echo "Creating sysfs for CNODE { COM, LTE, MASK }"
		#COM
		cd ${SYSFSDIR}
		create_sysfs_for_module "cnode" 1 1
		#LTE
		cd ${SYSFSDIR}
		create_sysfs_for_module "cnode" 2 0
		#MASK  
		cd ${SYSFSDIR}
		create_sysfs_for_module "cnode" 3 0
		;;
	"anode")				
		echo "Creating sysfs for ANODE { RF-CTRL , RF-FE }"
		cd ${SYSFSDIR}
		#RF-FE BOARD
		create_sysfs_for_module "anode" 4 0
		cd ${SYSFSDIR}
		# RF-CTRL-BOARD
		create_sysfs_for_module  "anode" 5 1
		;;
	*)
		echo "Wrong Unit Type."
		;;
esac
sync;
}


mkdir -p ${SYSDIR} ${SYSFSDIRHWMON} ${SYSFSDIRGPIO} ${SYSFSLED}
cd ${SYSDIRHWMON}
CURPWD=`pwd`

if [  $# -lt 1 ]; then
	usage;
	exit 0;
fi

while [[ $# -gt 0 ]]
do
	arg="$1"
	case "${arg}" in
	-u| --unittype)
		UNITIN="$2"	
       		create_sysfs_for_unit ${CURPWD} ${UNITIN}
		shift 
		shift
		;;
	-c| --clean)
		clean_sysfs_dir
		echo "SYSFS directory cleaned."
		shift
		exit 0
		;;
	*)
		echo "Unkown arguments."
		usage
		exit 0
		;;
	esac
done

echo "Script Complete"

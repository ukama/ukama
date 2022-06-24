#!/bin/bash
#Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

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
		ln -s ${EPATH}eeprom /tmp/sys/${UNITINFO}_inventory_db
	fi
}

create_and_write_to_file() {
	FILE=$1
	VALUE=$2

	touch ${FILE}
	[ -f "${FILE}" ] && { echo "${FILE} created."; }
	echo ${VALUE} > ${FILE}
	echo "Reading ${FILE}" `cat ${FILE}`
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

	create_and_write_to_file ${FTEMPVALUE} 45000
	create_and_write_to_file ${FMINVALUE}  25000
	create_and_write_to_file ${FMAXVALUE}  75000
	create_and_write_to_file ${FCRITVALUE} 85000
	create_and_write_to_file ${FCRITHYST}  2000
	create_and_write_to_file ${FMAXHYST}   2000
	create_and_write_to_file ${FOFFSET}    5000
	create_and_write_to_file ${FMINALARM}  0
	create_and_write_to_file ${FMAXALARM}  0
	create_and_write_to_file ${FCRITALARM} 0
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

	create_and_write_to_file in0_attvalue 63
	create_and_write_to_file in0_latch    0
}

create_sysfs_led() {
	SDIR=$1
	mkdir -p ${SDIR};
    cd ${SDIR};

	create_and_write_to_file brightness     0
	create_and_write_to_file max_brightness 255
	create_and_write_to_file trigger        "none"
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
    create_sysfs_led green
	create_sysfs_led blue	
}	


create_sysfs_for_ads1015() {
	SYSDIRPATH=$1
	NUM=$2
	SDIR=adc${NUM}
	cd ${SYSDIRPATH}
	mkdir -p ${SDIR}
	cd ${SDIR}

	create_and_write_to_file in0_input $RANDOM
	create_and_write_to_file in1_input $RANDOM
	create_and_write_to_file in2_input $RANDOM
	create_and_write_to_file in3_input $RANDOM
	create_and_write_to_file in4_input $RANDOM
	create_and_write_to_file in5_input $RANDOM
	create_and_write_to_file in6_input $RANDOM
	create_and_write_to_file in7_input $RANDOM
}

create_sysfs_for_inpgpio() {
	SYSDIRPATH=${SYSFSDIRGPIO}
	GPIONUM=$1
	cd ${SYSDIRPATH};
	mkdir -p gpio${GPIONUM}
	cd gpio${GPIONUM}

	create_and_write_to_file direction  "in"
	create_and_write_to_file value      1
	create_and_write_to_file edge       "rising"
	create_and_write_to_file active_low ""
	create_and_write_to_file polairy    0
}

create_sysfs_for_outgpio() {
	SYSDIRPATH=${SYSFSDIRGPIO}
    GPIONUM=$1
    cd ${SYSDIRPATH};
    mkdir -p gpio${GPIONUM}
    cd gpio${GPIONUM}

	create_and_write_to_file direction  "out"
	create_and_write_to_file value      1
	create_and_write_to_file edge       "both"
	create_and_write_to_file active_low ""
	create_and_write_to_file polairy    0
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
			create_sysfs_for_tmp464 ${MODDIR} 2
			create_sysfs_for_se98   ${MODDIR} 1
			create_sysfs_for_ina226 ${MODDIR} 1
			create_sysfs_for_ina226 ${MODDIR} 2
			create_sysfs_for_ina226 ${MODDIR} 3
            create_sysfs_for_leds   ${MODDIR} 0
            create_sysfs_for_leds   ${MODDIR} 1
            create_sysfs_for_leds   ${MODDIR} 2
            create_sysfs_for_leds   ${MODDIR} 3
            ;;

		2)
			create_sysfs_for_eeprom 1 0050 ${UNIT} ${MASTER}
            cd ${MODDIR}
			create_sysfs_for_tmp464 ${MODDIR} 1
			create_sysfs_for_tmp464 ${MODDIR} 2
			create_sysfs_for_se98   ${MODDIR} 1
			create_sysfs_for_ina226 ${MODDIR} 1
			create_sysfs_for_ina226 ${MODDIR} 2
			;;
		3)
			create_sysfs_for_eeprom 1 0051 ${UNIT} ${MASTER}
            cd ${MODDIR}
			create_sysfs_for_adt7481 ${MODDIR} 1
			create_sysfs_for_ina226  ${MODDIR} 1
			;;
		4)
			create_sysfs_for_eeprom 1 0052 ${UNIT} ${MASTER} 
			cd ${MODDIR}
            create_sysfs_for_se98    ${MODDIR} 1
			create_sysfs_for_tmp464  ${MODDIR} 1
			create_sysfs_for_ads1015 ${MODDIR} 1
			create_sysfs_for_att     ${MODDIR} 1
			create_sysfs_for_att     ${MODDIR} 2
			create_sysfs_for_inpgpio 38
            create_sysfs_for_inpgpio 35
            create_sysfs_for_inpgpio 34
            create_sysfs_for_outgpio 63
            create_sysfs_for_outgpio 61
            create_sysfs_for_outgpio 40
			;;
		5)
			create_sysfs_for_eeprom 0 0051 ${UNIT} ${MASTER}
            cd ${MODDIR}
			create_sysfs_for_se98   ${MODDIR} 1
            create_sysfs_for_tmp464 ${MODDIR} 1
			create_sysfs_for_leds   ${MODDIR} 0
            create_sysfs_for_leds   ${MODDIR} 1
            create_sysfs_for_leds   ${MODDIR} 2
            create_sysfs_for_leds   ${MODDIR} 3
			;;
		*)
			echo "Unknown module number."
			;;
		esac
}

usage() {
	echo "./prepare_env.sh [option]"
	echo "Options:"
	echo " -c | --clean		Clean the sysfs dir"
	echo " -u | --unittype        Create sysfs for unit (hnode, tnode, anode)"
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
	"hnode")
		echo "Creating sysfs for homeNode {Modules: TRX}"
		# TRX
		cd ${SYSFSDIR}
		create_sysfs_for_module "hnode" 2 1
		;;
	"tnode")
		echo "Creating sysfs for toweNode {Modules: COM, TRX, MASK}"
		# COM
		cd ${SYSFSDIR}
		create_sysfs_for_module "tnode" 1 1
		# TRX
		cd ${SYSFSDIR}
		create_sysfs_for_module "tnode" 2 0
		# MASK
		cd ${SYSFSDIR}
		create_sysfs_for_module "tnode" 3 0
		;;
	"anode")				
		echo "Creating sysfs for amplifierNode {Modules: CTRL, FE}"
		cd ${SYSFSDIR}
		# FE
		create_sysfs_for_module "anode" 4 0
		cd ${SYSFSDIR}
		# CTRL
		create_sysfs_for_module "anode" 5 1
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

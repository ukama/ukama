/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef ERRCODE_H_
#define ERRCODE_H_

//extern int *__geterrno(void);
//#define errno (*__geterrno())


#define ERR_UBSP_WR_FAIL 					    (-1001)
#define ERR_UBSP_R_FAIL  					    (-1002)
#define ERR_UBSP_INVALID_JSON_OBJECT 			(-1003)
#define ERR_UBSP_MW_ERR 					    (-1004)
#define ERR_UBSP_DB_MISSING_FIELD			    (-1005)
#define ERR_UBSP_INVALID_FIELD			        (-1006)
#define ERR_UBSP_DISABLED_FIELD			        (-1007)
#define ERR_UBSP_CRC_FAILURE				    (-1008)
#define ERR_UBSP_VALIDATION_FAILURE             (-1009)
#define ERR_UBSP_DB_MISSING_INFO				(-1010)
#define ERR_UBSP_JSON_PARSER                    (-1011)
#define ERR_UBSP_DB_MISSING_UNIT_INFO		    (-1012)
#define ERR_UBSP_INVALID_UNIT_INFO              (-1013)
#define ERR_UBSP_DB_MISSING_UNIT_CFG            (-1014)
#define ERR_UBSP_INVALID_UNIT_CFG               (-1015)
#define ERR_UBSP_DB_MISSING_MODULE_INFO	    	(-1016)
#define ERR_UBSP_INVALID_MODULE_INFO            (-1017)
#define ERR_UBSP_DB_MISSING_MODULE_CFG		    (-1018)
#define ERR_UBSP_INVALID_MODULE_CFG             (-1019)
#define ERR_UBSP_DB_MISSING_DEVICE_CFG		    (-1020)
#define ERR_UBSP_INVALID_DEVICE_CFG             (-1021)

#define ERR_UBSP_DESERIAL_FAIL                  (-1050)
#define ERR_UBSP_DB_MISSING_UNIT		    	(-1051)
#define ERR_UBSP_INVALID_UNIT                   (-1052)
#define ERR_UBSP_DB_MISSING_MODULE              (-1053)
#define ERR_UBSP_INVALID_MODULE                 (-1054)

#define ERR_UBSP_DEV_MISSING				    (-1101)
#define ERR_UBSP_DEV_PROPERTY_MISSING		    (-1102)
#define ERR_UBSP_DEV_HWATTR_MISSING		    	(-1103)
#define ERR_UBSP_DEV_DRVR_MISSING			    (-1104)
#define	ERR_UBSP_DEV_API_NOT_SUPPORTED			(-1105)
#define ERR_DRVR_API_NOT_SUPPORTED		        (-1115)
#define ERR_DRVR_API_NOT_AVAILABLE		        (-1116)

#define ERR_UBSP_SYSFS_FILE_MISSING			    (-1150)
#define ERR_UBSP_SYSFS_WRITE_FAILED			    (-1151)
#define ERR_UBSP_SYSFS_READ_FAILED				(-1152)

#define ERR_UBSP_DEV_IRQ_NOT_REG			    (-1161)

#define ERR_UBSP_THREAD_CREATE_FAIL			    (-1501)
#define ERR_UBSP_THREAD_CANCEL_FAIL			    (-1502)

#define ERR_UBSP_MEMORY_EXHAUSTED				(-1151)
#define ERR_UBSP_INVALID_POINTER				(-1151)

#define ERR_UBSP_LIST_DEL_FAILED				(-1601)

#define ERR_UBSP_UNEXPECTED_JSON_OBJECT         (-1701)
#define ERR_UBSP_CRT_JSON_SCHEMA         		(-1702)

#define ERR_UBSP_DB_LNK_MISSING                 (-1998)
#define ERR_UBSP_DB_MISSING                     (-1999)

#define ERR_EDGEREG_NOTAVAIL					(-1200)
#define ERR_EDGEREG_REGCREFAILED				(-1201)
#define ERR_EDGEREG_RESR_NOTIMPL				(-1202)
#define ERR_EDGEREG_INST_NOTAVAIL				(-1203)
#define ERR_EDGEREG_PROP_NOTAVAIL				(-1204)
#define ERR_EDGEREG_PROP_PERMDENIED				(-1205)
#define ERR_EDGEREG_PROP_NORESRAVAIL			(-1206)
#define ERR_EDGEREG_INST_INAVLIDOP				(-1207)
#define ERR_EDGEREG_INST_FXNTBL_MSNG			(-1208)

/* Socket */
#define ERR_SOCK_CREATION						(-1300)
#define ERR_SOCK_CONNECT						(-1301)
#define ERR_SOCK_SEND							(-1302)
#define ERR_SOCK_RECV							(-1303)
#define ERR_IFMSG_SERIALIZATION					(-1310)
#define ERR_IFMSG_DESERIALIZATION				(-1311)

#define ERR_IFMSG_MISMAT_TOKEN					(-1320)
#define ERR_IFMSG_MISMAT_INST					(-1320)
#define ERR_IFMSG_MISMAT_MSG_REQ				(-1321)
#define ERR_IFMSG_MISMAT_RSRC_ID				(-1322)
#define ERR_IFMSG_MISMAT_MSG_TYPE				(-1323)
#define ERR_UNEXPECTED_RESP_MSG					(-1324)
#endif /* ERRCODE_H_*/

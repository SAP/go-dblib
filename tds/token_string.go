// SPDX-FileCopyrightText: 2020-2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "stringer -type=Token"; DO NOT EDIT.

package tds

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TDS_CURDECLARE3-16]
	_ = x[TDS_PARAMFMT2-32]
	_ = x[TDS_LANGUAGE-33]
	_ = x[TDS_ORDERBY2-34]
	_ = x[TDS_CURDECLARE2-35]
	_ = x[TDS_COLFMTOLD-42]
	_ = x[TDS_DEBUGCMD-96]
	_ = x[TDS_ROWFMT2-97]
	_ = x[TDS_DYNAMIC2-98]
	_ = x[TDS_MSG-101]
	_ = x[TDS_LOGOUT-113]
	_ = x[TDS_OFFSET-120]
	_ = x[TDS_RETURNSTATUS-121]
	_ = x[TDS_PROCID-124]
	_ = x[TDS_CURCLOSE-128]
	_ = x[TDS_CURDELETE-129]
	_ = x[TDS_CURFETCH-130]
	_ = x[TDS_CURINFO-131]
	_ = x[TDS_CUROPEN-132]
	_ = x[TDS_CURUPDATE-133]
	_ = x[TDS_CURDECLARE-134]
	_ = x[TDS_CURINFO2-135]
	_ = x[TDS_CURINFO3-136]
	_ = x[TDS_COLNAME-160]
	_ = x[TDS_COLFMT-161]
	_ = x[TDS_EVENTNOTICE-162]
	_ = x[TDS_TABNAME-164]
	_ = x[TDS_COLINFO-165]
	_ = x[TDS_OPTIONCMD-166]
	_ = x[TDS_ALTNAME-167]
	_ = x[TDS_ALTFMT-168]
	_ = x[TDS_ORDERBY-169]
	_ = x[TDS_ERROR-170]
	_ = x[TDS_INFO-171]
	_ = x[TDS_RETURNVALUE-172]
	_ = x[TDS_LOGINACK-173]
	_ = x[TDS_CONTROL-174]
	_ = x[TDS_ALTCONTROL-175]
	_ = x[TDS_KEY-202]
	_ = x[TDS_ROW-209]
	_ = x[TDS_ALTROW-211]
	_ = x[TDS_PARAMS-215]
	_ = x[TDS_RPC-224]
	_ = x[TDS_CAPABILITY-226]
	_ = x[TDS_ENVCHANGE-227]
	_ = x[TDS_EED-229]
	_ = x[TDS_DBRPC-230]
	_ = x[TDS_DYNAMIC-231]
	_ = x[TDS_DBRPC2-232]
	_ = x[TDS_PARAMFMT-236]
	_ = x[TDS_ROWFMT-238]
	_ = x[TDS_DONE-253]
	_ = x[TDS_DONEPROC-254]
	_ = x[TDS_DONEINPROC-255]
}

const _Token_name = "TDS_CURDECLARE3TDS_PARAMFMT2TDS_LANGUAGETDS_ORDERBY2TDS_CURDECLARE2TDS_COLFMTOLDTDS_DEBUGCMDTDS_ROWFMT2TDS_DYNAMIC2TDS_MSGTDS_LOGOUTTDS_OFFSETTDS_RETURNSTATUSTDS_PROCIDTDS_CURCLOSETDS_CURDELETETDS_CURFETCHTDS_CURINFOTDS_CUROPENTDS_CURUPDATETDS_CURDECLARETDS_CURINFO2TDS_CURINFO3TDS_COLNAMETDS_COLFMTTDS_EVENTNOTICETDS_TABNAMETDS_COLINFOTDS_OPTIONCMDTDS_ALTNAMETDS_ALTFMTTDS_ORDERBYTDS_ERRORTDS_INFOTDS_RETURNVALUETDS_LOGINACKTDS_CONTROLTDS_ALTCONTROLTDS_KEYTDS_ROWTDS_ALTROWTDS_PARAMSTDS_RPCTDS_CAPABILITYTDS_ENVCHANGETDS_EEDTDS_DBRPCTDS_DYNAMICTDS_DBRPC2TDS_PARAMFMTTDS_ROWFMTTDS_DONETDS_DONEPROCTDS_DONEINPROC"

var _Token_map = map[Token]string{
	16:  _Token_name[0:15],
	32:  _Token_name[15:28],
	33:  _Token_name[28:40],
	34:  _Token_name[40:52],
	35:  _Token_name[52:67],
	42:  _Token_name[67:80],
	96:  _Token_name[80:92],
	97:  _Token_name[92:103],
	98:  _Token_name[103:115],
	101: _Token_name[115:122],
	113: _Token_name[122:132],
	120: _Token_name[132:142],
	121: _Token_name[142:158],
	124: _Token_name[158:168],
	128: _Token_name[168:180],
	129: _Token_name[180:193],
	130: _Token_name[193:205],
	131: _Token_name[205:216],
	132: _Token_name[216:227],
	133: _Token_name[227:240],
	134: _Token_name[240:254],
	135: _Token_name[254:266],
	136: _Token_name[266:278],
	160: _Token_name[278:289],
	161: _Token_name[289:299],
	162: _Token_name[299:314],
	164: _Token_name[314:325],
	165: _Token_name[325:336],
	166: _Token_name[336:349],
	167: _Token_name[349:360],
	168: _Token_name[360:370],
	169: _Token_name[370:381],
	170: _Token_name[381:390],
	171: _Token_name[390:398],
	172: _Token_name[398:413],
	173: _Token_name[413:425],
	174: _Token_name[425:436],
	175: _Token_name[436:450],
	202: _Token_name[450:457],
	209: _Token_name[457:464],
	211: _Token_name[464:474],
	215: _Token_name[474:484],
	224: _Token_name[484:491],
	226: _Token_name[491:505],
	227: _Token_name[505:518],
	229: _Token_name[518:525],
	230: _Token_name[525:534],
	231: _Token_name[534:545],
	232: _Token_name[545:555],
	236: _Token_name[555:567],
	238: _Token_name[567:577],
	253: _Token_name[577:585],
	254: _Token_name[585:597],
	255: _Token_name[597:611],
}

func (i Token) String() string {
	if str, ok := _Token_map[i]; ok {
		return str
	}
	return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
}

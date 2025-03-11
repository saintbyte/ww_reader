package ww_reader

// LILType Light intensity level
type LILType struct {
	Lumen int16
}

// Lightning flash detection
// "LFD,N,H,0"  //Уровень шума слишком высок
// "LFD,D,D,0"   //Disturber detected
// "LFD,L,D,%d"   //Lightning detected , and distanse
type LFDType struct {
	Type     string
	Value    string
	Distance int16
}

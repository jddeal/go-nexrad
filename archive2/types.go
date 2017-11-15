package archive2

const (
	RadialStatusStartOfElevationScan   = 0
	RadialStatusIntermediateRadialData = 1
	RadialStatusEndOfElevation         = 2
	RadialStatusBeginningOfVolumeScan  = 3
	RadialStatusEndOfVolumeScan        = 4
	RadialStatusStartNewElevation      = 5
)

// VolumeHeaderRecord for NEXRAD Archive II Data Streams
//
// Description:
// The Volume Header Record
// The Volume Header Record is fixed length and contains information uniquely
// identifying the format and the data that follows. Sits at the beginning of the
// Archive II data stream.
//
// Volume Header Record Data Format:
// The first 9 bytes is a character constant of which the last 2 characters
// identify the version. The next 3 bytes is a numeric string field starting
// with the value 001 and increasing by one for each volume of radar data in the
// queue to a maximum value of 999. Once the maximum value is reached the value
// will be rolled over. The combined 12 bytes are called the Archive II filename.
// The next 4 bytes contain the NEXRAD-modified Julian date the volume was
// produced at the RDA followed by 4 bytes containing the time the volume was
// recorded. The date and time integer values are big Endian. The last 4 bytes
// contain a 4-letter radar identifier assigned by ICAO.
//
// Version Number Reference:
// Version 02: Super Resolution disabled at the RDA (pre RDA Build 12.0)
// Version 03: Super Resolution (pre RDA Build 12.0)
// Version 04: Recombined Super Resolution
// Version 05: Super Resolution disabled at the RDA (RDA Build 12.0 and later)
// Version 06: Super Resolution (RDA Build 12.0 and later)
// Version 07: Recombined Super Resolution (RDA Build 12.0 and later)
type VolumeHeaderRecord struct {
	Tape      [7]byte
	Version   [2]byte
	Extension [3]byte
	// Date NEXRAD- modified Julian
	Date [4]byte
	// Time ms since midnight
	Time [4]byte
	// ICAO Radar identifier in ASCII. The four uppercase character International Civil Aviation Organization identifier of the radar producing the data.
	ICAO [4]byte
}

// LDMRecord (Local Data Manager) contains NEXRAD message data.
// Following the Volume Header Record are variable-length records containing the
// Archive II data messages. These records are referred to as LDM Compressed Record(s).
type LDMRecord struct {
	Size           int32
	MetaDataRecord []byte
}

type MessageHeader struct {
	MessageSize         uint16
	RDARedundantChannel uint8
	MessageType         uint8
	IDSequenceNumber    uint16
	JulianDate          uint16
	MillisOfDay         uint32
	NumMessageSegments  uint16
	MessageSegmentNum   uint16
}

// Message31 Digital Radar Data Generic Format
//
// Description:
// The message consists of base data information, that is, reflectivity, mean
// radial velocity, spectrum width, differential reflectivity, differential
// phase, correlation coefficient, azimuth angle, elevation angle, cut type,
// scanning strategy and calibration parameters. The frequency and volume of the
// message will be dependent on the scanning strategy and the type of data
// associated with that scanning strategy.
type Message31 struct {
	Header     Message31Header
	MomentData []interface{}
}

// RadialStatus enumerations
// 00: Start of new Elevation
// 01: Intermediate Radial Data
// 02: End of Elevation
// 03: Beginning of Volume Scan
// 04: End of Volume Scan
// 05: Start of new Elevation - Last Elevation in VCP

type RDAStatusMessage2 struct {
	RDAStatus                       uint16
	OperabilityStatus               uint16
	ControlStatus                   uint16
	AuxPowerGeneratorState          uint16
	AvgTxPower                      uint16
	HorizRefCalibCorr               uint16
	DataTxEnabled                   uint16
	VolumeCoveragePatternNum        uint16
	RDAControlAuth                  uint16
	RDABuild                        uint16
	OperationalMode                 uint16
	SuperResStatus                  uint16
	ClutterMitigationDecisionStatus uint16
	AvsetStatus                     uint16
	RDAAlarmSummary                 uint16
	CommandAck                      uint16
	ChannelControlStatus            uint16
	SpotBlankingStatus              uint16
	BypassMapGenDate                uint16
	BypassMapGenTime                uint16
	ClutterFilterMapGenDate         uint16
	ClutterFilterMapGenTime         uint16
	VertRefCalibCorr                uint16
	TransitionPwrSourceStatus       uint16
	RMSControlStatus                uint16
	PerformanceCheckStatus          uint16
	AlarmCodes                      uint16
}

type Message31Header struct {
	RadarIdentifier [4]byte
	// CollectionTime Radial data collection time in milliseconds past midnight GMT
	CollectionTime uint32
	// ModifiedJulianDate Current Julian date - 2440586.5
	ModifiedJulianDate uint16
	// AzimuthNumber Radial number within elevation scan
	AzimuthNumber uint16
	// AzimuthAngle Azimuth angle at which radial data was collected
	AzimuthAngle float32
	// CompressionIndicator Indicates if message type 31 is compressed and what method of compression is used. The Data Header Block is not compressed.
	CompressionIndicator uint8
	Spare                uint8
	// RadialLength Uncompressed length of the radial in bytes including the Data Header block length
	RadialLength uint16
	// AzimuthResolutionSpacing Azimuthal spacing between adjacent radials
	AzimuthResolutionSpacing uint8
	// RadialStatus Radial Status
	RadialStatus uint8
	// ElevationNumber Elevation number within volume scan
	ElevationNumber uint8
	// CutSectorNumber Sector Number within cut
	CutSectorNumber uint8
	// ElevationAngle Elevation angle at which radial radar data was collected
	ElevationAngle float32
	// RadialSpotBlankingStatus Spot blanking status for current radial, elevation scan and volume scan
	RadialSpotBlankingStatus uint8
	// AzimuthIndexingMode Azimuth indexing value (Set if azimuth angle is keyed to constant angles)
	AzimuthIndexingMode uint8
	DataBlockCount      uint16
	DataBlockPointers   [9]uint32
}

type DataBlock struct {
	DataBlockType [1]byte
	DataName      [3]byte
}

type VolumeData struct {
	DataBlock
	// LRTUP Size of data block in bytes
	LRTUP                          uint16
	VersionMajor                   uint8
	VersionMinor                   uint8
	Lat                            float32
	Long                           float32
	SiteHeight                     uint16
	FeedhornHeight                 uint16
	CalibrationConstant            float32
	SHVTXPowerHor                  float32
	SHVTXPowerVer                  float32
	SystemDifferentialReflectivity float32
	InitialSystemDifferentialPhase float32
	VolumeCoveragePatternNumber    uint16
	ProcessingStatus               uint16
}

type ElevationData struct {
	DataBlock
	// LRTUP Size of data block in bytes
	LRTUP uint16
	// ATMOS Atmospheric Attenuation Factor
	ATMOS [2]byte
	// CalibConst Scaling constant used by the Signal Processor for this elevation to calculate reflectivity
	CalibConst float32
}

type RadialData struct {
	DataBlock
	// LRTUP Size of data block in bytes
	LRTUP uint16
	// UnambiguousRange, Interval Size
	UnambiguousRange   uint16
	NoiseLevelHorz     float32
	NoiseLevelVert     float32
	NyquistVelocity    uint16
	Spares             [2]byte
	CalibConstHorzChan float32
	CalibConstVertChan float32
}

type GenericDataMoment struct {
	DataBlock
	Reserved uint32
	// NumberDataMomentGates Number of data moment gates for current radial
	NumberDataMomentGates uint16
	// DataMomentRange Range to center of first range gate
	DataMomentRange uint16
	// DataMomentRangeSampleInterval Size of data moment sample interval
	DataMomentRangeSampleInterval uint16
	// TOVER Threshold parameter which specifies the minimum difference in echo power between two resolution gates for them not to be labeled "overlayed"
	TOVER uint16
	// SNRThreshold SNR threshold for valid data
	SNRThreshold uint16
	// ControlFlags Indicates special control features
	ControlFlags uint8
	// DataWordSize Number of bits (DWS) used for storing data for each Data Moment gate
	DataWordSize uint8
	// Scale value used to convert Data Moments from integer to floating point data
	Scale float32
	// Offset value used to convert Data Moments from integer to floating point data
	Offset float32
}

type DataMoment struct {
	GenericDataMoment
	Data []byte
}

func (d *DataMoment) RefData() []float32 {
	convertedData := []float32{}
	for i := range d.Data {
		convertedData = append(convertedData, ScaleUint(d.Data[i], d.GenericDataMoment.Offset, d.GenericDataMoment.Scale))
	}
	return convertedData
}

// ScaleUint converts unsigned integer data that can be converted to floating point
// data using the Scale and Offset fields, i.e., F = (N - OFFSET) / SCALE where
// N is the integer data value and F is the resulting floating point value. A
// scale value of 0 indicates floating point moment data for each range gate.
func ScaleUint(n uint8, offset, scale float32) float32 {
	return (float32(n) - offset) / scale
}
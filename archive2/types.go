package archive2

import "time"

const (
	radialStatusStartOfElevationScan   = 0
	radialStatusIntermediateRadialData = 1
	radialStatusEndOfElevation         = 2
	radialStatusBeginningOfVolumeScan  = 3
	radialStatusEndOfVolumeScan        = 4
	radialStatusStartNewElevation      = 5

	// LegacyCTMHeaderLength sits in front of every message header
	LegacyCTMHeaderLength = 12

	// DefaultMetadataRecordLength is the size of every record regardless of its contents
	DefaultMetadataRecordLength = 2432
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
	Filename     [12]byte
	ModifiedDate int32   // ModifiedJulianDate NEXRAD date since 1970/1/1 = 1
	ModifiedTime int32   // X_ModifiedTime  ms since midnight
	ICAO         [4]byte // ICAO Radar identifier in ASCII. The four uppercase character International Civil Aviation Organization identifier of the radar producing the data.
}

// Date returns a time type representing the date of the scan capture
func (vh VolumeHeaderRecord) Date() time.Time {
	return timeFromModifiedJulian(int(vh.ModifiedDate), int(vh.ModifiedTime))
}

// FileName returns the name of the File
func (vh VolumeHeaderRecord) FileName() string {
	return string(vh.Filename[:])
}

func timeFromModifiedJulian(days, ms int) time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, int(days)).
		Add(time.Duration(ms) * time.Millisecond)
}

// LDMRecord (Local Data Manager) contains NEXRAD message data.
// Following the Volume Header Record are variable-length records containing the
// Archive II data messages. These records are referred to as LDM Compressed Record(s).
type LDMRecord struct {
	Size           int32
	MetaDataRecord []byte
}

// MessageHeader wrapper for archive2 Message Headers
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
	Header           Message31Header
	VolumeData       VolumeData
	ElevationData    ElevationData
	RadialData       RadialData
	ReflectivityData *DataMoment
	VelocityData     *DataMoment
	SwData           *DataMoment
	ZdrData          *DataMoment
	PhiData          *DataMoment
	RhoData          *DataMoment
	CfpData          *DataMoment
}

// Message2 contains RDA Status information about the radar site.
type Message2 struct {
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
	Spares                          [20]byte
}

// Message31Header contains header information for an Archive 2 Message 31 type
type Message31Header struct {
	RadarIdentifier              [4]byte
	CollectionTime               uint32  // CollectionTime Radial data collection time in milliseconds past midnight GMT
	ModifiedJulianDate           uint16  // ModifiedJulianDate Current Julian date - 2440586.5
	AzimuthNumber                uint16  // AzimuthNumber Radial number within elevation scan
	AzimuthAngle                 float32 // AzimuthAngle Azimuth angle at which radial data was collected
	CompressionIndicator         uint8   // CompressionIndicator Indicates if message type 31 is compressed and what method of compression is used. The Data Header Block is not compressed.
	Spare                        uint8
	RadialLength                 uint16  // RadialLength Uncompressed length of the radial in bytes including the Data Header block length
	AzimuthResolutionSpacingCode uint8   // AzimuthResolutionSpacing Code for the Azimuthal spacing between adjacent radials. 1 = .5 degrees, 2 = 1degree
	RadialStatus                 uint8   // RadialStatus Radial Status
	ElevationNumber              uint8   // ElevationNumber Elevation number within volume scan
	CutSectorNumber              uint8   // CutSectorNumber Sector Number within cut
	ElevationAngle               float32 // ElevationAngle Elevation angle at which radial radar data was collected
	RadialSpotBlankingStatus     uint8   // RadialSpotBlankingStatus Spot blanking status for current radial, elevation scan and volume scan
	AzimuthIndexingMode          uint8   // AzimuthIndexingMode Azimuth indexing value (Set if azimuth angle is keyed to constant angles)
	DataBlockCount               uint16
	DataBlockPointers            [10]uint32
}

// AzimuthResolutionSpacing returns the spacing in degrees according to the AzimuthResolutionSpacingCode
func (h *Message31Header) AzimuthResolutionSpacing() float64 {
	if h.AzimuthResolutionSpacingCode == 1 {
		return 0.5
	}
	return 1
}

// DataBlock wraps Data Block information
type DataBlock struct {
	DataBlockType [1]byte
	DataName      [3]byte
}

// VolumeData wraps information about the Volume being extracted
type VolumeData struct {
	LRTUP                          uint16 // LRTUP Size of data block in bytes
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

// ElevationData wraps Message 31 elevation data
type ElevationData struct {
	LRTUP      uint16  // LRTUP Size of data block in bytes
	ATMOS      [2]byte // ATMOS Atmospheric Attenuation Factor
	CalibConst float32 // CalibConst Scaling constant used by the Signal Processor for this elevation to calculate reflectivity
}

// RadialData wraps Message 31 radial data
type RadialData struct {
	LRTUP              uint16 // LRTUP Size of data block in bytes
	UnambiguousRange   uint16 // UnambiguousRange, Interval Size
	NoiseLevelHorz     float32
	NoiseLevelVert     float32
	NyquistVelocity    uint16
	Spares             [2]byte
	CalibConstHorzChan float32
	CalibConstVertChan float32
}

// GenericDataMoment is a generic data wrapper for momentary data. ex: REF, VEL, SW data
type GenericDataMoment struct {
	Reserved                      uint32
	NumberDataMomentGates         uint16  // NumberDataMomentGates Number of data moment gates for current radial
	DataMomentRange               uint16  // DataMomentRange Range to center of first range gate
	DataMomentRangeSampleInterval uint16  // DataMomentRangeSampleInterval Size of data moment sample interval
	TOVER                         uint16  // TOVER Threshold parameter which specifies the minimum difference in echo power between two resolution gates for them not to be labeled "overlayed"
	SNRThreshold                  uint16  // SNRThreshold SNR threshold for valid data
	ControlFlags                  uint8   // ControlFlags Indicates special control features
	DataWordSize                  uint8   // DataWordSize Number of bits (DWS) used for storing data for each Data Moment gate
	Scale                         float32 // Scale value used to convert Data Moments from integer to floating point data
	Offset                        float32 // Offset value used to convert Data Moments from integer to floating point data
}

// DataMoment wraps all Momentary data records. ex: REF, VEL, SW data
type DataMoment struct {
	GenericDataMoment
	Data []byte
}

// MomentDataBelowThreshold ...
const MomentDataBelowThreshold = 999

// MomentDataFolded ...
const MomentDataFolded = 998

// ScaledData automatically scales the nexrad moment values to their actual values.
// For all data moment integer values N = 0 indicates received signal is below
// threshold and N = 1 indicates range folded data. Actual data range is N = 2
// through 255, or 1023 for data resolution size 8, and 10 bits respectively.
func (d *DataMoment) ScaledData() []float32 {
	scaledData := []float32{}
	for _, v := range d.Data {
		if v == 0 {
			// below threshold
			scaledData = append(scaledData, MomentDataBelowThreshold)
		} else if v == 1 {
			// range folded
			scaledData = append(scaledData, MomentDataFolded)
		} else {
			scaledData = append(scaledData, scaleUint(uint16(v), d.GenericDataMoment.Offset, d.GenericDataMoment.Scale))
		}
	}
	return scaledData
}

// scaleUint converts unsigned integer data that can be converted to floating point
// data using the Scale and Offset fields, i.e., F = (N - OFFSET) / SCALE where
// N is the integer data value and F is the resulting floating point value. A
// scale value of 0 indicates floating point moment data for each range gate.
func scaleUint(n uint16, offset, scale float32) float32 {
	if scale == 0 {
		return float32(n)
	}
	return (float32(n) - offset) / scale
}

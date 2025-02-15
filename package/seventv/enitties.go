package seventv

type Role struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position int32  `json:"position"`
	Color    int32  `json:"color"`
	Allowed  int64  `json:"allowed"`
	Denied   int64  `json:"denied"`
}

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
}

type File struct {
	Name       string `json:"name"`
	StaticName string `json:"static_name"`
	Width      int32  `json:"width"`
	Height     int32  `json:"height"`
	FrameCount int32  `json:"frame_count"`
	Size       int32  `json:"size"`
	Format     string `json:"format"`
}

type FileFormat int

func (f FileFormat) String() string {
	return FileFormatMap[f]
}

const (
	GIF FileFormat = iota
	PNG
	WEBP
	AVIF
	JPG
)

var FileFormatMap = map[FileFormat]string{
	GIF:  "GIF",
	PNG:  "PNG",
	WEBP: "WEBP",
	AVIF: "AVIF",
	JPG:  "JPG",
}

var FileFormatReverseMap = map[string]FileFormat{
	"GIF":  GIF,
	"PNG":  PNG,
	"WEBP": WEBP,
	"AVIF": AVIF,
	"JPG":  JPG,
}

type FileScaler int

func (f FileScaler) String() string {
	return FileScalerMap[f]
}

const (
	X1 FileScaler = iota
	X2
	X3
	X4
)

var FileScalerMap = map[FileScaler]string{
	X1: "1x",
	X2: "2x",
	X3: "3x",
	X4: "4x",
}

var FileScalerReverseMap = map[string]FileScaler{
	"1x": X1,
	"2x": X2,
	"3x": X3,
	"4x": X4,
}

type Host struct {
	URL   string `json:"url"`
	Files []File
}

type Emote struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Owner *User    `json:"owner"`
	State []string `json:"state"`
	Tags  []string `json:"tags"`
	Host  *Host    `json:"host"`
}

type GQLResponse struct {
	Data struct {
		Emotes struct {
			Items []GQLEmote `json:"items"`
		} `json:"emotes"`
	} `json:"data"`
}

type GQLEmote struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner struct {
		Username string `json:"username"`
	} `json:"owner"`
	Host struct {
		URL string `json:"url"`
	} `json:"host"`
}

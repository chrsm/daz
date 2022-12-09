package daz

import (
	"bytes"
	"encoding/json"
	"time"
)

// Dufs are DSON, usually compressed, and DSON is just a fancy JSON, so it is
// parseable by any json lib in Go.
type Duf struct {
	FileVersion string `json:"file_version"`

	AssetInfo AssetInfo `json:"asset_info"`

	GeometryLibrary []Geometry `json:"geometry_library"`
	NodeLibrary     []Node     `json:"node_library"`
	ModifierLibrary []Modifier `json:"modifier_library"`
	ImageLibrary    []Image    `json:"image_library"`
	MaterialLibrary []Material `json:"material_library"`

	Scene Scene `json:"scene"`
}

type Scene struct {
	Presentation  Presentation       `json:"presentation"`
	Nodes         []NodeInstance     `json:"nodes"`
	UVs           []UVSetInstance    `json:"uvs"`
	Modifiers     []ModifierInstance `json:"modifiers"`
	Materials     []MaterialInstance `json:"materials"`
	Animations    []ChannelAnimation `json:"animations"`
	CurrentCamera string             `json:"current_camera"`

	Extra []map[string]interface{} `json:"extra"`
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/uv_set_instance/start
type UVSetInstance struct {
	ID     string
	URL    string
	Parent string
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/modifier_instance/start
type ModifierInstance struct {
	ID     string
	Parent string
	URL    string
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/material_instance/start
type MaterialInstance struct {
	ID       string
	URL      string
	Geometry string
	Groups   []string
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/channel_animation/start
type ChannelAnimation struct {
	URL  string
	Keys []TimeValuePair
}

// TODO
type TimeValuePair struct {
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/presentation/start
type Presentation struct {
	// content-type path
	Type        string         `json:"type"`
	Label       string         `json:"label"`
	Description string         `json:"description"`
	IconLarge   string         `json:"icon_large"`
	IconSmall   string         `json:"icon_small"`
	Colors      [2]Float3Array `json:"colors"`
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/node_instance/start
type NodeInstance struct {
	ID            string
	URL           string
	Parent        string
	ParentInPlace string
	ConformTarget string
	Geometries    []GeometryInstance
	Preview       Preview
}

type GeometryInstance struct {
	ID  string
	URL string
}

type Preview struct {
	OrientedBox OrientedBox
	CenterPoint [3]float32
	EndPoint    [3]float32
	// XYZ, YZX, ZYX, ZXY, XZY, YXZ
	RotationOrder string
}

type OrientedBox struct {
	Min [3]float32
	Max [3]float32
}

type DazTime struct {
	time.Time
}

func (dt *DazTime) UnmarshalJSON(b []byte) error {
	theTime, err := time.Parse(time.RFC3339, string(b))
	if err != nil {
		return err
	}

	*dt = DazTime{theTime}

	return nil
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/asset_info/start
type AssetInfo struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	Contributor Contributor `json:"contributor"`

	Revision string  `json:"revision"`
	Modified DazTime `json:"modified"`
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/contributor/start
type Contributor struct {
	Author  string `json:"author"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/geometry/start
type Geometry struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	// polygon_mesh | subdivision_surface
	Type   string `json:"type"`
	Source string `json:"source"`

	EdgeInterpolationMode string        `json:"edge_interpolation_mode"`
	Vertices              []Float3Array `json:"vertices"`

	PolygonGroups         []string `json:"polygon_groups"`
	PolygonMaterialGroups []string `json:"polygon_material_groups"`
	Polylist              Polylist `json:"polylist"`

	DefaultUVSet string `json:"default_uv_set"`

	RootRegion Region   `json:"root_region"`
	Graft      Graft    `json:"graft"`
	Rigidity   Rigidity `json:"rigidity"`

	// "An array of objects that represent additional application-specific information for this object."
	Extra []map[string]interface{} `json:"extra"`
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/rigidity/start
type Rigidity struct {
	Weights FloatIndexedArray `json:"weights"`
	Groups  []RigidityGroup   `json:"groups"`
}

type RigidityGroup struct {
	ID           string `json:"id"`
	RotationMode string `json:"rotation_mode"`
	// none, primary, secondary, tertiary
	ScaleModes []string `json:"scale_modes"`

	ReferenceVertices IntArray `json:"reference_vertices"`
	MaskVertices      IntArray `json:"mask_vertices"`

	Reference      string   `json:"reference"`
	TransformNodes []string `json:"transform_nodes"`
}

// Graft defines how to graft one figure's geometry to another.
// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/graft/start
type Graft struct {
	VertexCount int       `json:"vertex_count"`
	PolyCount   int       `json:"poly_count"`
	VertexPairs Int2Array `json:"vertex_pairs"`
	HiddenPolys IntArray  `json:"hidden_polys"`
}

type Int2Array struct {
	Count  int      `json:"count"`
	Values [][2]int `json:"values"`
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/region/start
type Region struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	// cards_on | cards_off
	DisplayHint string   `json:"display_hint"`
	Map         IntArray `json:"map"`
	Children    []Region `json:"children"`
}

// http://docs.daz3d.com/doku.php/public/dson_spec/object_definitions/polygon/start
type Polylist struct {
	Count  int     `json:"count"`
	Values [][]int `json:"values"` // maxL=6
}

// http://docs.daz3d.com/doku.php/public/dson_spec/format_description/data_types/start#int_array
type IntArray struct {
	Count  int   `json:"count"`
	Values []int `json:"values"`
}

// http://docs.daz3d.com/doku.php/public/dson_spec/format_description/data_types/start#float3_array
type Float3Array struct {
	Count  int          `json:"count"`
	Values [][3]float32 `json:"values"`
}

type FloatIndexedArray struct {
	Count int `json:"count"`
	// pairs = [[int, float], ...]
	Values []FloatPair `json:"values"`
}

type floatidxiface struct {
	Count  int `json:"count"`
	Values [][2]interface{}
}

func (f3 *FloatIndexedArray) UnmarshalJSON(b []byte) error {

	dec := json.NewDecoder(bytes.NewReader(b))
	temp := &floatidxiface{}
	if err := dec.Decode(temp); err != nil {
		return err
	}

	f3.Count = temp.Count
	f3.Values = make([]FloatPair, len(temp.Values))
	for i := range temp.Values {
		f3.Values[i].Index = temp.Values[i][0].(int)
		f3.Values[i].Value = temp.Values[i][1].(float32)
	}

	return nil
}

type FloatPair struct {
	Index int
	Value float32
}

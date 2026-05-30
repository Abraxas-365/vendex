package kernel

// Commerce domain IDs

type ProductID string

func NewProductID(id string) ProductID { return ProductID(id) }
func (p ProductID) String() string     { return string(p) }
func (p ProductID) IsEmpty() bool      { return string(p) == "" }

type OrderID string

func NewOrderID(id string) OrderID { return OrderID(id) }
func (o OrderID) String() string   { return string(o) }
func (o OrderID) IsEmpty() bool    { return string(o) == "" }

type OrderItemID string

func NewOrderItemID(id string) OrderItemID { return OrderItemID(id) }
func (o OrderItemID) String() string       { return string(o) }
func (o OrderItemID) IsEmpty() bool        { return string(o) == "" }

type CustomerID string

func NewCustomerID(id string) CustomerID { return CustomerID(id) }
func (c CustomerID) String() string      { return string(c) }
func (c CustomerID) IsEmpty() bool       { return string(c) == "" }

type CategoryID string

func NewCategoryID(id string) CategoryID { return CategoryID(id) }
func (c CategoryID) String() string      { return string(c) }
func (c CategoryID) IsEmpty() bool       { return string(c) == "" }

type CollectionID string

func NewCollectionID(id string) CollectionID { return CollectionID(id) }
func (c CollectionID) String() string        { return string(c) }
func (c CollectionID) IsEmpty() bool         { return string(c) == "" }

type PageID string

func NewPageID(id string) PageID { return PageID(id) }
func (p PageID) String() string  { return string(p) }
func (p PageID) IsEmpty() bool   { return string(p) == "" }

type PageVersionID string

func NewPageVersionID(id string) PageVersionID { return PageVersionID(id) }
func (p PageVersionID) String() string         { return string(p) }
func (p PageVersionID) IsEmpty() bool          { return string(p) == "" }

type PromoID string

func NewPromoID(id string) PromoID { return PromoID(id) }
func (p PromoID) String() string   { return string(p) }
func (p PromoID) IsEmpty() bool    { return string(p) == "" }

type MediaID string

func NewMediaID(id string) MediaID { return MediaID(id) }
func (m MediaID) String() string   { return string(m) }
func (m MediaID) IsEmpty() bool    { return string(m) == "" }

type PluginID string

func NewPluginID(id string) PluginID { return PluginID(id) }
func (p PluginID) String() string    { return string(p) }
func (p PluginID) IsEmpty() bool     { return string(p) == "" }

type PluginVersionID string

func NewPluginVersionID(id string) PluginVersionID { return PluginVersionID(id) }
func (p PluginVersionID) String() string           { return string(p) }
func (p PluginVersionID) IsEmpty() bool            { return string(p) == "" }

type InstallationID string

func NewInstallationID(id string) InstallationID { return InstallationID(id) }
func (i InstallationID) String() string          { return string(i) }
func (i InstallationID) IsEmpty() bool           { return string(i) == "" }

type SettingID string

func NewSettingID(id string) SettingID { return SettingID(id) }
func (s SettingID) String() string     { return string(s) }
func (s SettingID) IsEmpty() bool      { return string(s) == "" }

type VendorID string

func NewVendorID(id string) VendorID { return VendorID(id) }
func (v VendorID) String() string    { return string(v) }
func (v VendorID) IsEmpty() bool     { return string(v) == "" }

type ThemeID string

func NewThemeID(id string) ThemeID { return ThemeID(id) }
func (t ThemeID) String() string   { return string(t) }
func (t ThemeID) IsEmpty() bool    { return string(t) == "" }

type BlockID string

func NewBlockID(id string) BlockID { return BlockID(id) }
func (b BlockID) String() string   { return string(b) }
func (b BlockID) IsEmpty() bool    { return string(b) == "" }

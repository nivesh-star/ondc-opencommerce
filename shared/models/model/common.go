package model

import "time"

// Every API call in beckn protocol has a context. It provides a high-level overview to the receiver about the nature of the intended transaction. Typically, it is the BAP that sets the transaction context based on the consumer's location and action on their UI. But sometimes, during unsolicited callbacks, the BPP also sets the transaction context but it is usually the same as the context of a previous full-cycle, request-callback interaction between the BAP and the BPP. The context object contains four types of fields. <ol><li>Demographic information about the transaction using fields like `domain`, `country`, and `region`.</li><li>Addressing details like the sending and receiving platform's ID and API URL.</li><li>Interoperability information like the protocol version that implemented by the sender and,</li><li>Transaction details like the method being called at the receiver's endpoint, the transaction_id that represents an end-to-end user session at the BAP, a message ID to pair requests with callbacks, a timestamp to capture sending times, a ttl to specifiy the validity of the request, and a key to encrypt information if necessary.</li></ol> This object must be passed in every interaction between a BAP and a BPP. In HTTP/S implementations, it is not necessary to send the context during the synchronous response. However, in asynchronous protocols, the context must be sent during all interactions,
type Context struct {
	// Domain code that is relevant to this transaction context
	Domain string `json:"domain,omitempty"`
	// The location where the transaction is intended to be fulfilled.
	Location *Location `json:"location,omitempty"`
	// The Beckn protocol method being called by the sender and executed at the receiver.
	Action string `json:"action,omitempty"`
	// Version of transaction protocol being used by the sender.
	Version string `json:"version,omitempty"`
	// Subscriber ID of the BAP
	BapId string `json:"bap_id,omitempty"`
	// Subscriber URL of the BAP for accepting callbacks from BPPs.
	BapUri string `json:"bap_uri,omitempty"`
	// Subscriber ID of the BPP
	BppId string `json:"bpp_id,omitempty"`
	// Subscriber URL of the BPP for accepting calls from BAPs.
	BppUri string `json:"bpp_uri,omitempty"`
	// This is a unique value which persists across all API calls from `search` through `confirm`. This is done to indicate an active user session across multiple requests. The BPPs can use this value to push personalized recommendations, and dynamic offerings related to an ongoing transaction despite being unaware of the user active on the BAP.
	TransactionId string `json:"transaction_id,omitempty"`
	// This is a unique value which persists during a request / callback cycle. Since beckn protocol APIs are asynchronous, BAPs need a common value to match an incoming callback from a BPP to an earlier call. This value can also be used to ignore duplicate messages coming from the BPP. It is recommended to generate a fresh message_id for every new interaction. When sending unsolicited callbacks, BPPs must generate a new message_id.
	MessageId string `json:"message_id,omitempty"`
	// Time of request generation in RFC3339 format
	Timestamp time.Time `json:"timestamp,omitempty"`
	// The encryption public key of the sender
	Key string `json:"key,omitempty"`
	// The duration in ISO8601 format after timestamp for which this message holds valid
	Ttl string `json:"ttl,omitempty"`
}

type Location struct {
	Id         string      `json:"id,omitempty"`
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	// The url to the map of the location. This can be a globally recognized map url or the one specified by the network policy.
	MapUrl string `json:"map_url,omitempty"`
	// The GPS co-ordinates of this location.
	Gps string `json:"gps,omitempty"`
	// The address of this location.
	Address string `json:"address,omitempty"`
	// The city this location is, or is located within
	City *City `json:"city,omitempty"`
	// The state this location is, or is located within
	District string `json:"district,omitempty"`
	// The state this location is, or is located within
	State *State `json:"state,omitempty"`
	// The country this location is, or is located within
	Country  *Country `json:"country,omitempty"`
	AreaCode string   `json:"area_code,omitempty"`
	Circle   *Circle  `json:"circle,omitempty"`
	// The boundary polygon of this location
	Polygon string `json:"polygon,omitempty"`
	// The three dimensional region describing this location
	Var3dspace string `json:"3dspace,omitempty"`
	// The rating of this location
	Rating string `json:"rating,omitempty"`
}

// Physical description of something.
type Descriptor struct {
	Name           string                    `json:"name,omitempty"`
	Code           string                    `json:"code,omitempty"`
	ShortDesc      string                    `json:"short_desc,omitempty"`
	LongDesc       string                    `json:"long_desc,omitempty"`
	AdditionalDesc *DescriptorAdditionalDesc `json:"additional_desc,omitempty"`
	Media          []MediaFile               `json:"media,omitempty"`
	Images         []Image                   `json:"images,omitempty"`
}

type DescriptorAdditionalDesc struct {
	Url         string `json:"url,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}

// This object contains a url to a media file.
type MediaFile struct {
	// indicates the nature and format of the document, file, or assortment of bytes. MIME types are defined and standardized in IETF's RFC 6838
	Mimetype string `json:"mimetype,omitempty"`
	// The URL of the file
	Url string `json:"url,omitempty"`
	// The digital signature of the file signed by the sender
	Signature string `json:"signature,omitempty"`
	// The signing algorithm used by the sender
	Dsa string `json:"dsa,omitempty"`
}

// Describes an image
type Image struct {
	// URL to the image. This can be a data url or an remote url
	Url string `json:"url,omitempty"`
	// The size of the image. The network policy can define the default dimensions of each type
	SizeType string `json:"size_type,omitempty"`
	// Width of the image in pixels
	Width string `json:"width,omitempty"`
	// Height of the image in pixels
	Height string `json:"height,omitempty"`
}
type City struct {
	// Name of the city
	Name string `json:"name,omitempty"`
	// City code
	Code string `json:"code,omitempty"`
}

// A bounded geopolitical region of governance inside a country.
type State struct {
	// Name of the state
	Name string `json:"name,omitempty"`
	// State code as per country or international standards
	Code string `json:"code,omitempty"`
}

// Describes a country
type Country struct {
	// Name of the country
	Name string `json:"name,omitempty"`
	// Country code as per ISO 3166-1 and ISO 3166-2 format
	Code string `json:"code,omitempty"`
}

// Describes a circular region of a specified radius centered at a specified GPS coordinate.
type Circle struct {
	Gps    string  `json:"gps,omitempty"`
	Radius *Scalar `json:"radius,omitempty"`
}

// Describes a scalar
type Scalar struct {
	Type_          string       `json:"type,omitempty"`
	Value          string       `json:"value,omitempty"`
	EstimatedValue string       `json:"estimated_value,omitempty"`
	ComputedValue  string       `json:"computed_value,omitempty"`
	Range_         *ScalarRange `json:"range,omitempty"`
	Unit           string       `json:"unit,omitempty"`
}

type ScalarRange struct {
	Min string `json:"min,omitempty"`
	Max string `json:"max,omitempty"`
}

// The intent to buy or avail a product or a service. The BAP can declare the intent of the consumer containing <ul><li>What they want (A product, service, offer)</li><li>Who they want (A seller, service provider, agent etc)</li><li>Where they want it and where they want it from</li><li>When they want it (start and end time of fulfillment</li><li>How they want to pay for it</li></ul><br>This has properties like descriptor,provider,fulfillment,payment,category,offer,item,tags<br>This is typically used by the BAP to send the purpose of the user's search to the BPP. This will be used by the BPP to find products or services it offers that may match the user's intent.<br>For example, in Mobility, the mobility consumer declares a mobility intent. In this case, the mobility consumer declares information that describes various aspects of their journey like,<ul><li>Where would they like to begin their journey (intent.fulfillment.start.location)</li><li>Where would they like to end their journey (intent.fulfillment.end.location)</li><li>When would they like to begin their journey (intent.fulfillment.start.time)</li><li>When would they like to end their journey (intent.fulfillment.end.time)</li><li>Who is the transport service provider they would like to avail services from (intent.provider)</li><li>Who is traveling (This is not recommended in public networks) (intent.fulfillment.customer)</li><li>What kind of fare product would they like to purchase (intent.item)</li><li>What add-on services would they like to avail</li><li>What offers would they like to apply on their booking (intent.offer)</li><li>What category of services would they like to avail (intent.category)</li><li>What additional luggage are they carrying</li><li>How would they like to pay for their journey (intent.payment)</li></ul><br>For example, in health domain, a consumer declares the intent for a lab booking the describes various aspects of their booking like,<ul><li>Where would they like to get their scan/test done (intent.fulfillment.start.location)</li><li>When would they like to get their scan/test done (intent.fulfillment.start.time)</li><li>When would they like to get the results of their test/scan (intent.fulfillment.end.time)</li><li>Who is the service provider they would like to avail services from (intent.provider)</li><li>Who is getting the test/scan (intent.fulfillment.customer)</li><li>What kind of test/scan would they like to purchase (intent.item)</li><li>What category of services would they like to avail (intent.category)</li><li>How would they like to pay for their journey (intent.payment)</li></ul>
type Intent struct {
	// A raw description of the search intent. Free text search strings, raw audio, etc can be sent in this object.
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	// The provider from which the customer wants to place to the order from
	Provider *Provider `json:"provider,omitempty"`
	// Details on how the customer wants their order fulfilled
	Fulfillment *Fulfillment `json:"fulfillment,omitempty"`
	// Details on how the customer wants to pay for the order
	Payment *Payment `json:"payment,omitempty"`
	// Details on the item category
	Category *Category `json:"category,omitempty"`
	// details on the offer the customer wants to avail
	Offer *Offer `json:"offer,omitempty"`
	// Details of the item that the consumer wants to order
	Item *Item      `json:"item,omitempty"`
	Tags []TagGroup `json:"tags,omitempty"`
}

// Describes the catalog of a business.
type Provider struct {
	// Id of the provider
	Id         string      `json:"id,omitempty"`
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	// Category Id of the provider at the BPP-level catalog
	CategoryId   string        `json:"category_id,omitempty"`
	Rating       string        `json:"rating,omitempty"`
	Time         *Time         `json:"time,omitempty"`
	Categories   []Category    `json:"categories,omitempty"`
	Fulfillments []Fulfillment `json:"fulfillments,omitempty"`
	Payments     []Payment     `json:"payments,omitempty"`
	Locations    []Location    `json:"locations,omitempty"`
	Offers       []Offer       `json:"offers,omitempty"`
	Items        []Item        `json:"items,omitempty"`
	// Time after which catalog has to be refreshed
	Exp time.Time `json:"exp,omitempty"`
	// Whether this provider can be rated or not
	Rateable bool `json:"rateable,omitempty"`
	// The time-to-live in seconds, for this object. This can be overriden at deeper levels. A value of -1 indicates that this object is not cacheable.
	Ttl  string     `json:"ttl,omitempty"`
	Tags []TagGroup `json:"tags,omitempty"`
}

// Describes time in its various forms. It can be a single point in time; duration; or a structured timetable of operations<br>This has properties like label, time stamp,duration,range, days, schedule
type Time struct {
	Label     string     `json:"label,omitempty"`
	Timestamp time.Time  `json:"timestamp,omitempty"`
	Duration  string     `json:"duration,omitempty"`
	Range_    *TimeRange `json:"range,omitempty"`
	// comma separated values representing days of the week
	Days     string    `json:"days,omitempty"`
	Schedule *Schedule `json:"schedule,omitempty"`
}
type TimeRange struct {
	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`
}

// Describes schedule as a repeating time period used to describe a regularly recurring event. At a minimum a schedule will specify frequency which describes the interval between occurrences of the event. Additional information can be provided to specify the schedule more precisely. This includes identifying the timestamps(s) of when the event will take place. Schedules may also have holidays to exclude a specific day from the schedule.<br>This has properties like frequency, holidays, times
type Schedule struct {
	Frequency string      `json:"frequency,omitempty"`
	Holidays  []time.Time `json:"holidays,omitempty"`
	Times     []time.Time `json:"times,omitempty"`
}

// A label under which a collection of items can be grouped.
type Category struct {
	// ID of the category
	Id               string      `json:"id,omitempty"`
	ParentCategoryId string      `json:"parent_category_id,omitempty"`
	Descriptor       *Descriptor `json:"descriptor,omitempty"`
	Time             *Time       `json:"time,omitempty"`
	// Time to live for an instance of this schema
	Ttl  string     `json:"ttl,omitempty"`
	Tags []TagGroup `json:"tags,omitempty"`
}

// A collection of tag objects with group level attributes. For detailed documentation on the Tags and Tag Groups schema go to https://github.com/beckn/protocol-specifications/discussions/316
type TagGroup struct {
	// Indicates the display properties of the tag group. If display is set to false, then the group will not be displayed. If it is set to true, it should be displayed. However, group-level display properties can be overriden by individual tag-level display property. As this schema is purely for catalog display purposes, it is not recommended to send this value during search.
	Display bool `json:"display,omitempty"`
	// Description of the TagGroup, can be used to store detailed information.
	Descriptor *TagDescriptor `json:"descriptor,omitempty"`
	// An array of Tag objects listed under this group. This property can be set by BAPs during search to narrow the `search` and achieve more relevant results. When received during `on_search`, BAPs must render this list under the heading described by the `name` property of this schema.
	List []Tag `json:"list,omitempty"`
}

// Description of the Tag, can be used to store detailed information.
type TagDescriptor struct {
	Name           string                    `json:"name,omitempty"`
	Code           string                    `json:"code,omitempty"`
	ShortDesc      string                    `json:"short_desc,omitempty"`
	LongDesc       string                    `json:"long_desc,omitempty"`
	AdditionalDesc *DescriptorAdditionalDesc `json:"additional_desc,omitempty"`
	Media          []MediaFile               `json:"media,omitempty"`
	Images         []Image                   `json:"images,omitempty"`
}

// Describes a tag. This is used to contain extended metadata. This object can be added as a property to any schema to describe extended attributes. For BAPs, tags can be sent during search to optimize and filter search results. BPPs can use tags to index their catalog to allow better search functionality. Tags are sent by the BPP as part of the catalog response in the `on_search` callback. Tags are also meant for display purposes. Upon receiving a tag, BAPs are meant to render them as name-value pairs. This is particularly useful when rendering tabular information about a product or service.
type Tag struct {
	// Description of the Tag, can be used to store detailed information.
	Descriptor *TagDescriptor `json:"descriptor,omitempty"`
	// The value of the tag. This set by the BPP and rendered as-is by the BAP.
	Value string `json:"value,omitempty"`
	// This value indicates if the tag is intended for display purposes. If set to `true`, then this tag must be displayed. If it is set to `false`, it should not be displayed. This value can override the group display value.
	Display bool `json:"display,omitempty"`
}

// Describes how a an order will be rendered/fulfilled to the end-customer
type Fulfillment struct {
	// Unique reference ID to the fulfillment of an order
	Id string `json:"id,omitempty"`
	// A code that describes the mode of fulfillment. This is typically set when there are multiple ways an order can be fulfilled. For example, a retail order can be fulfilled either via store pickup or a home delivery. Similarly, a medical consultation can be provided either in-person or via tele-consultation. The network policy must publish standard fulfillment type codes for the different modes of fulfillment.
	Type_ string `json:"type,omitempty"`
	// Whether the fulfillment can be rated or not
	Rateable bool `json:"rateable,omitempty"`
	// The rating value of the fulfullment service.
	Rating string `json:"rating,omitempty"`
	// The current state of fulfillment. The BPP must set this value whenever the state of the order fulfillment changes and fire an unsolicited `on_status` call.
	State *FulfillmentState `json:"state,omitempty"`
	// Indicates whether the fulfillment allows tracking
	Tracking bool `json:"tracking,omitempty"`
	// The person that will ultimately receive the order
	Customer *FulfillmentCustomer `json:"customer,omitempty"`
	// The agent that is currently handling the fulfillment of the order
	Agent   *Agent   `json:"agent,omitempty"`
	Contact *Contact `json:"contact,omitempty"`
	Vehicle *Vehicle `json:"vehicle,omitempty"`
	// The list of logical stops encountered during the fulfillment of an order.
	Stops []Stop `json:"stops,omitempty"`
	// The physical path taken by the agent that can be rendered on a map. The allowed format of this property can be set by the network.
	Path string     `json:"path,omitempty"`
	Tags []TagGroup `json:"tags,omitempty"`
}

// The current state of fulfillment. The BPP must set this value whenever the state of the order fulfillment changes and fire an unsolicited `on_status` call.
type FulfillmentState struct {
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	UpdatedAt  time.Time   `json:"updated_at,omitempty"`
	// ID of entity which changed the state
	UpdatedBy string `json:"updated_by,omitempty"`
}

// The person that will ultimately receive the order
type FulfillmentCustomer struct {
	Organization *Organization `json:"organization,omitempty"`
	Person       *Person       `json:"person,omitempty"`
	Contact      *Contact      `json:"contact,omitempty"`
}

// An organization. Usually a recognized business entity.
type Organization struct {
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	// The postal address of the organization
	Address string `json:"address,omitempty"`
	// The state where the organization's address is registered
	State *State `json:"state,omitempty"`
	// The city where the the organization's address is registered
	City    *City        `json:"city,omitempty"`
	Contact *Contact     `json:"contact,omitempty"`
	Creds   []Credential `json:"creds,omitempty"`
}

// Describes a person as any individual
type Person struct {
	// Describes the identity of the person
	Id string `json:"id,omitempty"`
	// Profile url of the person
	Url string `json:"url,omitempty"`
	// the name of the person
	Name  string `json:"name,omitempty"`
	Image *Image `json:"image,omitempty"`
	// Age of the person
	Age string `json:"age,omitempty"`
	// Date of birth of the person
	Dob string `json:"dob,omitempty"`
	// Gender of something, typically a Person, but possibly also fictional characters, animals, etc. While Male and Female may be used, text strings are also acceptable for people who do not identify as a binary gender.Allowed values for this field can be published in the network policy
	Gender    string            `json:"gender,omitempty"`
	Creds     []Credential      `json:"creds,omitempty"`
	Languages []PersonLanguages `json:"languages,omitempty"`
	Skills    []PersonSkills    `json:"skills,omitempty"`
	Tags      []TagGroup        `json:"tags,omitempty"`
}

// Describes the contact information of an entity
type Contact struct {
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	// A Jcard object as per draft-ietf-jcardcal-jcard-03 specification
	Jcard *interface{} `json:"jcard,omitempty"`
}

// Describes a credential of an entity - Person or Organization
type Credential struct {
	Id    string `json:"id,omitempty"`
	Type_ string `json:"type,omitempty"`
	// URL of the credential
	Url string `json:"url,omitempty"`
}

// Describes a language known to the person.
type PersonLanguages struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

// Describes a skill of the person.
type PersonSkills struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

// Describes the direct performer, driver or executor that fulfills an order. It is usually a person. But in some rare cases, it could be a non-living entity like a drone, or a bot. Some examples of agents are Doctor in the healthcare sector, a driver in the mobility sector, or a delivery person in the logistics sector. This object can be set at any stage of the order lifecycle. This can be set at the discovery stage when the BPP wants to provide details on the agent fulfilling the order, like in healthcare, where the doctor's name appears during search. This object can also used to search for a particular person that the customer wants fulfilling an order. Sometimes, this object gets instantiated after the order is confirmed, like in the case of on-demand taxis, where the driver is assigned after the user confirms the ride.
type Agent struct {
	Person       *Person       `json:"person,omitempty"`
	Contact      *Contact      `json:"contact,omitempty"`
	Organization *Organization `json:"organization,omitempty"`
	Rating       string        `json:"rating,omitempty"`
}
type Vehicle struct {
	Category         string `json:"category,omitempty"`
	Capacity         int32  `json:"capacity,omitempty"`
	Make             string `json:"make,omitempty"`
	Model            string `json:"model,omitempty"`
	Size             string `json:"size,omitempty"`
	Variant          string `json:"variant,omitempty"`
	Color            string `json:"color,omitempty"`
	EnergyType       string `json:"energy_type,omitempty"`
	Registration     string `json:"registration,omitempty"`
	WheelsCount      string `json:"wheels_count,omitempty"`
	CargoVolumne     string `json:"cargo_volumne,omitempty"`
	WheelchairAccess string `json:"wheelchair_access,omitempty"`
	Code             string `json:"code,omitempty"`
	EmissionStandard string `json:"emission_standard,omitempty"`
}

// A logical point in space and time during the fulfillment of an order.
type Stop struct {
	Id           string `json:"id,omitempty"`
	ParentStopId string `json:"parent_stop_id,omitempty"`
	// Location of the stop
	Location *Location `json:"location,omitempty"`
	// The type of stop. Allowed values of this property can be defined by the network policy.
	Type_ string `json:"type,omitempty"`
	// Timings applicable at the stop.
	Time *Time `json:"time,omitempty"`
	// Instructions that need to be followed at the stop
	Instructions *StopInstructions `json:"instructions,omitempty"`
	// Contact details of the stop
	Contact *Contact `json:"contact,omitempty"`
	// The details of the person present at the stop
	Person        *Person        `json:"person,omitempty"`
	Authorization *Authorization `json:"authorization,omitempty"`
}

// Instructions that need to be followed at the stop
type StopInstructions struct {
	Name           string                    `json:"name,omitempty"`
	Code           string                    `json:"code,omitempty"`
	ShortDesc      string                    `json:"short_desc,omitempty"`
	LongDesc       string                    `json:"long_desc,omitempty"`
	AdditionalDesc *DescriptorAdditionalDesc `json:"additional_desc,omitempty"`
	Media          []MediaFile               `json:"media,omitempty"`
	Images         []Image                   `json:"images,omitempty"`
}

// Describes an authorization mechanism used to start or end the fulfillment of an order. For example, in the mobility sector, the driver may require a one-time password to initiate the ride. In the healthcare sector, a patient may need to provide a password to open a video conference link during a teleconsultation.
type Authorization struct {
	// Type of authorization mechanism used. The allowed values for this field can be published as part of the network policy.
	Type_ string `json:"type,omitempty"`
	// Token used for authorization. This is typically generated at the BPP. The BAP can send this value to the user via any channel that it uses to authenticate the user like SMS, Email, Push notification, or in-app rendering.
	Token string `json:"token,omitempty"`
	// Timestamp in RFC3339 format from which token is valid
	ValidFrom time.Time `json:"valid_from,omitempty"`
	// Timestamp in RFC3339 format until which token is valid
	ValidTo time.Time `json:"valid_to,omitempty"`
	// Status of the token
	Status string `json:"status,omitempty"`
}

// Describes the terms of settlement between the BAP and the BPP for a single transaction. When instantiated, this object contains <ol><li>the amount that has to be settled,</li><li>The payment destination destination details</li><li>When the settlement should happen, and</li><li>A transaction reference ID</li></ol>. During a transaction, the BPP reserves the right to decide the terms of payment. However, the BAP can send its terms to the BPP first. If the BPP does not agree to those terms, it must overwrite the terms and return them to the BAP. If overridden, the BAP must either agree to the terms sent by the BPP in order to preserve the provider's autonomy, or abort the transaction. In case of such disagreements, the BAP and the BPP can perform offline negotiations on the payment terms. Once an agreement is reached, the BAP and BPP can resume transactions.
type Payment struct {
	// ID of the payment term that can be referred at an item or an order level in a catalog
	Id string `json:"id,omitempty"`
	// This field indicates who is the collector of payment. The BAP can set this value to 'bap' if it wants to collect the payment first and  settle it to the BPP. If the BPP agrees to those terms, the BPP should not send the payment url. Alternatively, the BPP can set this field with the value 'bpp' if it wants the payment to be made directly.
	CollectedBy string `json:"collected_by,omitempty"`
	// A payment url to be called by the BAP. If empty, then the payment is to be done offline. The details of payment should be present in the params object. If tl_method = http/get, then the payment details will be sent as url params. Two url param values, ```$transaction_id``` and ```$amount``` are mandatory.
	Url    string         `json:"url,omitempty"`
	Params *PaymentParams `json:"params,omitempty"`
	Type_  string         `json:"type,omitempty"`
	Status string         `json:"status,omitempty"`
	Time   *Time          `json:"time,omitempty"`
	Tags   []TagGroup     `json:"tags,omitempty"`
}
type PaymentParams struct {
	// The reference transaction ID associated with a payment activity
	TransactionId               string `json:"transaction_id,omitempty"`
	Amount                      string `json:"amount,omitempty"`
	Currency                    string `json:"currency,omitempty"`
	BankCode                    string `json:"bank_code,omitempty"`
	BankAccountNumber           string `json:"bank_account_number,omitempty"`
	VirtualPaymentAddress       string `json:"virtual_payment_address,omitempty"`
	SourceBankCode              string `json:"source_bank_code,omitempty"`
	SourceBankAccountNumber     string `json:"source_bank_account_number,omitempty"`
	SourceVirtualPaymentAddress string `json:"source_virtual_payment_address,omitempty"`
}

// An offer associated with a catalog. This is typically used to promote a particular product and enable more purchases.
type Offer struct {
	Id          string      `json:"id,omitempty"`
	Descriptor  *Descriptor `json:"descriptor,omitempty"`
	LocationIds []string    `json:"location_ids,omitempty"`
	CategoryIds []string    `json:"category_ids,omitempty"`
	ItemIds     []string    `json:"item_ids,omitempty"`
	Time        *Time       `json:"time,omitempty"`
	Tags        []TagGroup  `json:"tags,omitempty"`
}

// Describes a product or a service offered to the end consumer by the provider. In the mobility sector, it can represent a fare product like one way journey. In the logistics sector, it can represent the delivery service offering. In the retail domain it can represent a product like a grocery item.
type Item struct {
	// ID of the item.
	Id string `json:"id,omitempty"`
	// ID of the item, this item is a variant of
	ParentItemId string `json:"parent_item_id,omitempty"`
	// The number of units of the parent item this item is a multiple of
	ParentItemQuantity *ItemQuantity `json:"parent_item_quantity,omitempty"`
	// Physical description of the item
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	// The creator of this item
	Creator *ItemCreator `json:"creator,omitempty"`
	// The price of this item, if it has intrinsic value
	Price *Price `json:"price,omitempty"`
	// The selling quantity of the item
	Quantity *ItemQuantity `json:"quantity,omitempty"`
	// Categories this item can be listed under
	CategoryIds []string `json:"category_ids,omitempty"`
	// Modes through which this item can be fulfilled
	FulfillmentIds []string `json:"fulfillment_ids,omitempty"`
	// Provider Locations this item is available in
	LocationIds []string `json:"location_ids,omitempty"`
	// Payment modalities through which this item can be ordered
	PaymentIds []string `json:"payment_ids,omitempty"`
	AddOns     []AddOn  `json:"add_ons,omitempty"`
	// Cancellation terms of this item
	CancellationTerms []CancellationTerm `json:"cancellation_terms,omitempty"`
	// Refund terms of this item
	RefundTerms []Terms `json:"refund_terms,omitempty"`
	// Terms that are applicable be met when this item is replaced
	ReplacementTerms []Terms `json:"replacement_terms,omitempty"`
	// Terms that are applicable when this item is returned
	ReturnTerms []Terms `json:"return_terms,omitempty"`
	// Additional input required from the customer to purchase / avail this item
	Xinput *XInput `json:"xinput,omitempty"`
	// Temporal attributes of this item. This property is used when the item exists on the catalog only for a limited period of time.
	Time *Time `json:"time,omitempty"`
	// Whether this item can be rated
	Rateable bool `json:"rateable,omitempty"`
	// The rating of the item
	Rating string `json:"rating,omitempty"`
	// Whether this item is an exact match of the request
	Matched bool `json:"matched,omitempty"`
	// Whether this item is a related item to the exactly matched item
	Related bool `json:"related,omitempty"`
	// Whether this item is a recommended item to a response
	Recommended bool `json:"recommended,omitempty"`
	// Time to live in seconds for an instance of this schema
	Ttl  string     `json:"ttl,omitempty"`
	Tags []TagGroup `json:"tags,omitempty"`
}

// Refund term of an item or an order
type Terms struct {
	// The state of fulfillment during which this term is applicable.
	FulfillmentState *TermsFulfillmentState `json:"fulfillment_state,omitempty"`
	// Indicates if cancellation will result in a refund
	RefundEligible bool `json:"refund_eligible,omitempty"`
	// Time within which refund will be processed after successful cancellation.
	RefundWithin *TermsRefundWithin `json:"refund_within,omitempty"`
	RefundAmount *Price             `json:"refund_amount,omitempty"`
}

// The state of fulfillment during which this term is applicable.
type TermsFulfillmentState struct {
	// Name of the state
	Name string `json:"name,omitempty"`
	// State code as per country or international standards
	Code string `json:"code,omitempty"`
}

// Time within which refund will be processed after successful cancellation.
type TermsRefundWithin struct {
	Label     string     `json:"label,omitempty"`
	Timestamp time.Time  `json:"timestamp,omitempty"`
	Duration  string     `json:"duration,omitempty"`
	Range_    *TimeRange `json:"range,omitempty"`
	// comma separated values representing days of the week
	Days     string    `json:"days,omitempty"`
	Schedule *Schedule `json:"schedule,omitempty"`
}

// Describes the cancellation terms of an item or an order. This can be referenced at an item or order level. Item-level cancellation terms can override the terms at the order level.
type CancellationTerm struct {
	// The state of fulfillment during which this term is applicable.
	FulfillmentState *FulfillmentState `json:"fulfillment_state,omitempty"`
	// Indicates whether a reason is required to cancel the order
	ReasonRequired bool `json:"reason_required,omitempty"`
	// Information related to the time of cancellation.
	CancelBy        *CancelBy  `json:"cancel_by,omitempty"`
	CancellationFee *Fee       `json:"cancellation_fee,omitempty"`
	Xinput          *XInput    `json:"xinput,omitempty"`
	ExternalRef     *MediaFile `json:"external_ref,omitempty"`
}

// Contains any additional or extended inputs required to confirm an order. This is typically a Form Input. Sometimes, selection of catalog elements is not enough for the BPP to confirm an order. For example, to confirm a flight ticket, the airline requires details of the passengers along with information on baggage, identity, in addition to the class of ticket. Similarly, a logistics company may require details on the nature of shipment in order to confirm the shipping. A recruiting firm may require additional details on the applicant in order to confirm a job application. For all such purposes, the BPP can choose to send this object attached to any object in the catalog that is required to be sent while placing the order. This object can typically be sent at an item level or at the order level. The item level XInput will override the Order level XInput as it indicates a special requirement of information for that particular item. Hence the BAP must render a separate form for the Item and another form at the Order level before confirmation.
type XInput struct {
	Head         *XInputHead         `json:"head,omitempty"`
	Form         *Form               `json:"form,omitempty"`
	FormResponse *XInputFormResponse `json:"form_response,omitempty"`
	// Indicates whether the form data is mandatorily required by the BPP to confirm the order.
	Required bool `json:"required,omitempty"`
}

// Describes a form
type Form struct {
	// The form identifier.
	Id string `json:"id,omitempty"`
	// The URL from where the form can be fetched. The content fetched from the url must be processed as per the mime_type specified in this object. Once fetched, the rendering platform can choosed to render the form as-is as an embeddable element; or process it further to blend with the theme of the application. In case the interface is non-visual, the the render can process the form data and reproduce it as per the standard specified in the form.
	Url string `json:"url,omitempty"`
	// The form submission data
	Data map[string]string `json:"data,omitempty"`
	// This field indicates the nature and format of the form received by querying the url. MIME types are defined and standardized in IETF's RFC 6838.
	MimeType            string `json:"mime_type,omitempty"`
	Resubmit            bool   `json:"resubmit,omitempty"`
	MultipleSumbissions bool   `json:"multiple_sumbissions,omitempty"`
}

// Describes the response to a form submission
type XInputFormResponse struct {
	// Contains the status of form submission.
	Status       string       `json:"status,omitempty"`
	Signature    string       `json:"signature,omitempty"`
	SubmissionId string       `json:"submission_id,omitempty"`
	Errors       []ModelError `json:"errors,omitempty"`
}

// Describes an error object that is returned by a BAP, BPP or BG as a response or callback to an action by another network participant. This object is sent when any request received by a network participant is unacceptable. This object can be sent either during Ack or with the callback.
type ModelError struct {
	// Standard error code. For full list of error codes, refer to docs/protocol-drafts/BECKN-005-ERROR-CODES-DRAFT-01.md of this repo\"
	Code string `json:"code,omitempty"`
	// Path to json schema generating the error. Used only during json schema validation errors
	Paths string `json:"paths,omitempty"`
	// Human readable message describing the error. Used mainly for logging. Not recommended to be shown to the user.
	Message string `json:"message,omitempty"`
}

// Provides the header information for the xinput.
type XInputHead struct {
	Descriptor *Descriptor      `json:"descriptor,omitempty"`
	Index      *XInputHeadIndex `json:"index,omitempty"`
	Headings   []string         `json:"headings,omitempty"`
}

type XInputHeadIndex struct {
	Min int32 `json:"min,omitempty"`
	Cur int32 `json:"cur,omitempty"`
	Max int32 `json:"max,omitempty"`
}

// A fee applied on a particular entity
type Fee struct {
	// Percentage of a value
	Percentage string `json:"percentage,omitempty"`
	// A fixed value
	Amount *AllOfFeeAmount `json:"amount,omitempty"`
}

// A fixed value
type AllOfFeeAmount struct {
	Currency       string `json:"currency,omitempty"`
	Value          string `json:"value,omitempty"`
	EstimatedValue string `json:"estimated_value,omitempty"`
	ComputedValue  string `json:"computed_value,omitempty"`
	ListedValue    string `json:"listed_value,omitempty"`
	OfferedValue   string `json:"offered_value,omitempty"`
	MinimumValue   string `json:"minimum_value,omitempty"`
	MaximumValue   string `json:"maximum_value,omitempty"`
}

// Information related to the time of cancellation.
type CancelBy struct {
	Label     string     `json:"label,omitempty"`
	Timestamp time.Time  `json:"timestamp,omitempty"`
	Duration  string     `json:"duration,omitempty"`
	Range_    *TimeRange `json:"range,omitempty"`
	// comma separated values representing days of the week
	Days     string    `json:"days,omitempty"`
	Schedule *Schedule `json:"schedule,omitempty"`
}

// Describes an additional item offered as a value-addition to a product or service. This does not exist independently in a catalog and is always associated with an item.
type AddOn struct {
	// Provider-defined ID of the add-on
	Id         string        `json:"id,omitempty"`
	Descriptor *Descriptor   `json:"descriptor,omitempty"`
	Price      *Price        `json:"price,omitempty"`
	Quantity   *ItemQuantity `json:"quantity,omitempty"`
}

// Describes the price of a product or service
type Price struct {
	Currency       string `json:"currency,omitempty"`
	Value          string `json:"value,omitempty"`
	EstimatedValue string `json:"estimated_value,omitempty"`
	ComputedValue  string `json:"computed_value,omitempty"`
	ListedValue    string `json:"listed_value,omitempty"`
	OfferedValue   string `json:"offered_value,omitempty"`
	MinimumValue   string `json:"minimum_value,omitempty"`
	MaximumValue   string `json:"maximum_value,omitempty"`
}

// The creator of this item
type ItemCreator struct {
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	// The postal address of the organization
	Address string `json:"address,omitempty"`
	// The state where the organization's address is registered
	State *State `json:"state,omitempty"`
	// The city where the the organization's address is registered
	City    *City        `json:"city,omitempty"`
	Contact *Contact     `json:"contact,omitempty"`
	Creds   []Credential `json:"creds,omitempty"`
}

// Describes the count or amount of an item
type ItemQuantity struct {
	Allocated *ItemQuantityAllocated `json:"allocated,omitempty"`
	Available *ItemQuantityAvailable `json:"available,omitempty"`
	Maximum   *ItemQuantityMaximum   `json:"maximum,omitempty"`
	Minimum   *ItemQuantityMinimum   `json:"minimum,omitempty"`
	Selected  *ItemQuantitySelected  `json:"selected,omitempty"`
	Unitized  *ItemQuantityUnitized  `json:"unitized,omitempty"`
}

// This represents the exact quantity allocated for purchase of the item.
type ItemQuantityAllocated struct {
	Count   int32   `json:"count,omitempty"`
	Measure *Scalar `json:"measure,omitempty"`
}

// This represents the exact quantity available for purchase of the item. The buyer can only purchase multiples of this
type ItemQuantityAvailable struct {
	Count   int32   `json:"count,omitempty"`
	Measure *Scalar `json:"measure,omitempty"`
}

// This represents the maximum quantity allowed for purchase of the item
type ItemQuantityMaximum struct {
	Count   int32   `json:"count,omitempty"`
	Measure *Scalar `json:"measure,omitempty"`
}

// This represents the quantity selected for purchase of the item
type ItemQuantitySelected struct {
	Count   int32   `json:"count,omitempty"`
	Measure *Scalar `json:"measure,omitempty"`
}

// This represents the minimum quantity allowed for purchase of the item
type ItemQuantityMinimum struct {
	Count   int32   `json:"count,omitempty"`
	Measure *Scalar `json:"measure,omitempty"`
}

// This represents the quantity available in a single unit of the item
type ItemQuantityUnitized struct {
	Count   int32   `json:"count,omitempty"`
	Measure *Scalar `json:"measure,omitempty"`
}

//***********************************************************************************************************************
// ********************************************* END ON SEARCH API DS  **************************************************
// ********************************************* START OF ON SELECT DS **************************************************

// Describes a legal purchase order. It contains the complete details of the legal contract created between the buyer and the seller.
type Order struct {
	// Human-readable ID of the order. This is generated at the BPP layer. The BPP can either generate order id within its system or forward the order ID created at the provider level.
	Id string `json:"id,omitempty"`
	// A list of order IDs to link this order to previous orders.
	RefOrderIds []string `json:"ref_order_ids,omitempty"`
	// Status of the order. Allowed values can be defined by the network policy
	Status string `json:"status,omitempty"`
	// This is used to indicate the type of order being created to BPPs. Sometimes orders can be linked to previous orders, like a replacement order in a retail domain. A follow-up consultation in healthcare domain. A single order part of a subscription order. The list of order types can be standardized at the network level.
	Type_ string `json:"type,omitempty"`
	// Details of the provider whose catalog items have been selected.
	Provider *Provider `json:"provider,omitempty"`
	// The items purchased / availed in this order
	Items []Item `json:"items,omitempty"`
	// The add-ons purchased / availed in this order
	AddOns []AddOn `json:"add_ons,omitempty"`
	// The offers applied in this order
	Offers []Offer `json:"offers,omitempty"`
	// The billing details of this order
	Billing *Billing `json:"billing,omitempty"`
	// The fulfillments involved in completing this order
	Fulfillments []Fulfillment `json:"fulfillments,omitempty"`
	// The cancellation details of this order
	Cancellation *Cancellation `json:"cancellation,omitempty"`
	// Cancellation terms of this item
	CancellationTerms []CancellationTerm `json:"cancellation_terms,omitempty"`
	Documents         []OrderDocuments   `json:"documents,omitempty"`
	// Refund terms of this item
	RefundTerms []Terms `json:"refund_terms,omitempty"`
	// Replacement terms of this item
	ReplacementTerms []Terms `json:"replacement_terms,omitempty"`
	// Return terms of this item
	ReturnTerms []Terms `json:"return_terms,omitempty"`
	// The mutually agreed upon quotation for this order.
	Quote *Quotation `json:"quote,omitempty"`
	// The terms of settlement for this order
	Payments []Payment `json:"payments,omitempty"`
	// The date-time of creation of this order
	CreatedAt time.Time `json:"created_at,omitempty"`
	// The date-time of updated of this order
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Additional input required from the customer to confirm this order
	Xinput *XInput    `json:"xinput,omitempty"`
	Tags   []TagGroup `json:"tags,omitempty"`
}

// Describes a quote. It is the estimated price of products or services from the BPP.<br>This has properties like price, breakup, ttl
type Quotation struct {
	// ID of the quote.
	Id string `json:"id,omitempty"`
	// The total quoted price
	Price *Price `json:"price,omitempty"`
	// the breakup of the total quoted price
	Breakup []QuotationBreakup `json:"breakup,omitempty"`
	Ttl     string             `json:"ttl,omitempty"`
}
type QuotationBreakup struct {
	Item  *Item  `json:"item,omitempty"`
	Title string `json:"title,omitempty"`
	Price *Price `json:"price,omitempty"`
}

// Documnents associated to the order
type OrderDocuments struct {
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	// This field indicates the nature and format of the form received by querying the url. MIME types are defined and standardized in IETF's RFC 6838.
	MimeType string `json:"mime_type,omitempty"`
	// The URL from where the form can be fetched. The content fetched from the url must be processed as per the mime_type specified in this object.
	Url string `json:"url,omitempty"`
	// The URL from where the form can be fetched. The content fetched from the url must be processed as per the mime_type specified in this object.
	OldPolicyDoc string `json:"old_policy_doc,omitempty"`
}

// Describes the billing details of an entity.<br>This has properties like name,organization,address,email,phone,time,tax_number, created_at,updated_at
type Billing struct {
	// Name of the billable entity
	Name string `json:"name,omitempty"`
	// Details of the organization being billed.
	Organization *Organization `json:"organization,omitempty"`
	// The address of the billable entity
	Address string `json:"address,omitempty"`
	// The state where the billable entity resides. This is important for state-level tax calculation
	State *State `json:"state,omitempty"`
	// The city where the billable entity resides.
	City *City `json:"city,omitempty"`
	// Email address where the bill is sent to
	Email string `json:"email,omitempty"`
	// Phone number of the billable entity
	Phone string `json:"phone,omitempty"`
	// Details regarding the billing period
	Time *Time `json:"time,omitempty"`
	// ID of the billable entity as recognized by the taxation authority
	TaxId string `json:"tax_id,omitempty"`
}

// Describes a cancellation event
type Cancellation struct {
	// Date-time when the order was cancelled by the buyer
	Time        time.Time `json:"time,omitempty"`
	CancelledBy string    `json:"cancelled_by,omitempty"`
	// The reason for cancellation
	Reason *CancellationReason `json:"reason,omitempty"`
	// Any additional information regarding the nature of cancellation
	AdditionalDescription *CancellationAdditionalDescription `json:"additional_description,omitempty"`
}

// The reason for cancellation
type CancellationReason struct {
	Id         string      `json:"id,omitempty"`
	Descriptor *Descriptor `json:"descriptor,omitempty"`
}
type CancellationAdditionalDescription struct {
	Name           string                    `json:"name,omitempty"`
	Code           string                    `json:"code,omitempty"`
	ShortDesc      string                    `json:"short_desc,omitempty"`
	LongDesc       string                    `json:"long_desc,omitempty"`
	AdditionalDesc *DescriptorAdditionalDesc `json:"additional_desc,omitempty"`
	Media          []MediaFile               `json:"media,omitempty"`
	Images         []Image                   `json:"images,omitempty"`
}

// ***********************************************************************************************************************
// ********************************************* END ON ON-SELECT API DS  ***********************************************
// ********************************************* START OF RATING  DS *************************************************
// Describes the rating of an entity
type Rating struct {
	// Category of the entity being rated
	RatingCategory string `json:"rating_category,omitempty"`
	// Id of the object being rated
	Id string `json:"id,omitempty"`
	// Rating value given to the object. This can be a single value or can also contain an inequality operator like gt, gte, lt, lte. This can also contain an inequality expression containing logical operators like && and ||.
	Value string `json:"value,omitempty"`
}

// ***********************************************************************************************************************
// ********************************************* END ON RATING API DS  ***********************************************
// ********************************************* START OF BAP  DS *************************************************

// Describes the products or services offered by a BPP. This is typically sent as the response to a search intent from a BAP. The payment terms, offers and terms of fulfillment supported by the BPP can also be included here. The BPP can show hierarchical nature of products/services in its catalog using the parent_category_id in categories. The BPP can also send a ttl (time to live) in the context which is the duration for which a BAP can cache the catalog and use the cached catalog.  <br>This has properties like bbp/descriptor,bbp/categories,bbp/fulfillments,bbp/payments,bbp/offers,bbp/providers and exp<br>This is used in the following situations.<br><ul><li>This is typically used in the discovery stage when the BPP sends the details of the products and services it offers as response to a search intent from the BAP. </li></ul>
type Catalog struct {
	Descriptor *Descriptor `json:"descriptor,omitempty"`
	// Fulfillment modes offered at the BPP level. This is used when a BPP itself offers fulfillments on behalf of the providers it has onboarded.
	Fulfillments []Fulfillment `json:"fulfillments,omitempty"`
	// Payment terms offered by the BPP for all transactions. This can be overriden at the provider level.
	Payments []Payment `json:"payments,omitempty"`
	// Offers at the BPP-level. This is common across all providers onboarded by the BPP.
	Offers    []Offer    `json:"offers,omitempty"`
	Providers []Provider `json:"providers,omitempty"`
	// Timestamp after which catalog will expire
	Exp time.Time `json:"exp,omitempty"`
	// Duration in seconds after which this catalog will expire
	Ttl  string     `json:"ttl,omitempty"`
	Tags []TagGroup `json:"tags,omitempty"`
}

// Contains tracking information that can be used by the BAP to track the fulfillment of an order in real-time. which is useful for knowing the location of time sensitive deliveries.
type Tracking struct {
	// A unique tracking reference number
	Id string `json:"id,omitempty"`
	// A URL to the tracking endpoint. This can be a link to a tracking webpage, a webhook URL created by the BAP where BPP can push the tracking data, or a GET url creaed by the BPP which the BAP can poll to get the tracking data. It can also be a websocket URL where the BPP can push real-time tracking data.
	Url string `json:"url,omitempty"`
	// In case there is no real-time tracking endpoint available, this field will contain the latest location of the entity being tracked. The BPP will update this value everytime the BAP calls the track API.
	Location *Location `json:"location,omitempty"`
	// This value indicates if the tracking is currently active or not. If this value is `active`, then the BAP can begin tracking the order. If this value is `inactive`, the tracking URL is considered to be expired and the BAP should stop tracking the order.
	Status string `json:"status,omitempty"`
}
type OnRatingBody struct {
	Context *Context         `json:"context"`
	Message *OnRatingMessage `json:"message"`
	Error_  *ModelError      `json:"error,omitempty"`
}

type OnRatingMessage struct {
	// A feedback form to allow the user to provide additional information on the rating provided
	FeedbackForm *MessageFeedbackForm `json:"feedback_form,omitempty"`
}

// A feedback form to allow the user to provide additional information on the rating provided
type MessageFeedbackForm struct {
	Head         *XInputHead         `json:"head,omitempty"`
	Form         *Form               `json:"form,omitempty"`
	FormResponse *XInputFormResponse `json:"form_response,omitempty"`
	// Indicates whether the form data is mandatorily required by the BPP to confirm the order.
	Required bool `json:"required,omitempty"`
}

// ***********************************************************************************************************************
// ********************************************* END of BAP API DS  ***********************************************
// *********************************************  DS used in code   *************************************************

// Describes the acknowledgement sent in response to an API call. If the implementation uses HTTP/S, then Ack must be returned in the same session. Every API call to a BPP must be responded to with an Ack whether the BPP intends to respond with a callback or not. This has one property called `status` that indicates the status of the Acknowledgement.
type Ack struct {
	// The status of the acknowledgement. If the request passes the validation criteria of the BPP, then this is set to ACK. If a BPP responds with status = `ACK` to a request, it is required to respond with a callback. If the request fails the validation criteria, then this is set to NACK. Additionally, if a BPP does not intend to respond with a callback even after the request meets the validation criteria, it should set this value to `NACK`.
	Status string `json:"status,omitempty"`
	// A list of tags containing any additional information sent along with the Acknowledgement.
	Tags []TagGroup `json:"tags,omitempty"`
}

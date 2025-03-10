// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.12.4
// source: recipe/recipe.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// GetRecipeRequest is used to request a specific recipe.
type GetRecipeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RecipeId      string                 `protobuf:"bytes,1,opt,name=recipe_id,json=recipeId,proto3" json:"recipe_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRecipeRequest) Reset() {
	*x = GetRecipeRequest{}
	mi := &file_recipe_recipe_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRecipeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecipeRequest) ProtoMessage() {}

func (x *GetRecipeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_recipe_recipe_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecipeRequest.ProtoReflect.Descriptor instead.
func (*GetRecipeRequest) Descriptor() ([]byte, []int) {
	return file_recipe_recipe_proto_rawDescGZIP(), []int{0}
}

func (x *GetRecipeRequest) GetRecipeId() string {
	if x != nil {
		return x.RecipeId
	}
	return ""
}

// GetRecipeResponse returns the full details of a recipe.
type GetRecipeResponse struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	RecipeId          string                 `protobuf:"bytes,1,opt,name=recipe_id,json=recipeId,proto3" json:"recipe_id,omitempty"`
	Title             string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Ingredients       string                 `protobuf:"bytes,3,opt,name=ingredients,proto3" json:"ingredients,omitempty"`
	Steps             string                 `protobuf:"bytes,4,opt,name=steps,proto3" json:"steps,omitempty"`
	NutritionalInfo   string                 `protobuf:"bytes,5,opt,name=nutritional_info,json=nutritionalInfo,proto3" json:"nutritional_info,omitempty"`
	AllergyDisclaimer string                 `protobuf:"bytes,6,opt,name=allergy_disclaimer,json=allergyDisclaimer,proto3" json:"allergy_disclaimer,omitempty"`
	Appliances        string                 `protobuf:"bytes,7,opt,name=appliances,proto3" json:"appliances,omitempty"`
	CreatedAt         int64                  `protobuf:"varint,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt         int64                  `protobuf:"varint,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *GetRecipeResponse) Reset() {
	*x = GetRecipeResponse{}
	mi := &file_recipe_recipe_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRecipeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecipeResponse) ProtoMessage() {}

func (x *GetRecipeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_recipe_recipe_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecipeResponse.ProtoReflect.Descriptor instead.
func (*GetRecipeResponse) Descriptor() ([]byte, []int) {
	return file_recipe_recipe_proto_rawDescGZIP(), []int{1}
}

func (x *GetRecipeResponse) GetRecipeId() string {
	if x != nil {
		return x.RecipeId
	}
	return ""
}

func (x *GetRecipeResponse) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *GetRecipeResponse) GetIngredients() string {
	if x != nil {
		return x.Ingredients
	}
	return ""
}

func (x *GetRecipeResponse) GetSteps() string {
	if x != nil {
		return x.Steps
	}
	return ""
}

func (x *GetRecipeResponse) GetNutritionalInfo() string {
	if x != nil {
		return x.NutritionalInfo
	}
	return ""
}

func (x *GetRecipeResponse) GetAllergyDisclaimer() string {
	if x != nil {
		return x.AllergyDisclaimer
	}
	return ""
}

func (x *GetRecipeResponse) GetAppliances() string {
	if x != nil {
		return x.Appliances
	}
	return ""
}

func (x *GetRecipeResponse) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *GetRecipeResponse) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

// RecipeQueryRequest is used for both advanced search and list operations.
// An empty "query" field indicates a listing operation, while a non-empty field
// triggers advanced search logic (e.g., filtering by cuisine, diet, etc.).
type RecipeQueryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Query         string                 `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`                 // Advanced search text (e.g., "vegan"); empty for simple listings.
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // Optional: Filter recipes by creator's user ID.
	Filter        string                 `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`               // Optional: Additional filtering criteria (e.g., "Indian").
	Page          int32                  `protobuf:"varint,4,opt,name=page,proto3" json:"page,omitempty"`                  // Optional: Requested page number for pagination.
	Limit         int32                  `protobuf:"varint,5,opt,name=limit,proto3" json:"limit,omitempty"`                // Optional: Number of recipes per page.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RecipeQueryRequest) Reset() {
	*x = RecipeQueryRequest{}
	mi := &file_recipe_recipe_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecipeQueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecipeQueryRequest) ProtoMessage() {}

func (x *RecipeQueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_recipe_recipe_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecipeQueryRequest.ProtoReflect.Descriptor instead.
func (*RecipeQueryRequest) Descriptor() ([]byte, []int) {
	return file_recipe_recipe_proto_rawDescGZIP(), []int{2}
}

func (x *RecipeQueryRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *RecipeQueryRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *RecipeQueryRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

func (x *RecipeQueryRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *RecipeQueryRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// RecipeQueryResponse returns the results for a query along with pagination details.
type RecipeQueryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Recipes       []*GetRecipeResponse   `protobuf:"bytes,1,rep,name=recipes,proto3" json:"recipes,omitempty"` // List of recipes matching the query.
	Page          int32                  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`      // Echoed page number.
	Limit         int32                  `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`    // Echoed limit per page.
	Total         int32                  `protobuf:"varint,4,opt,name=total,proto3" json:"total,omitempty"`    // Total number of matching recipes.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RecipeQueryResponse) Reset() {
	*x = RecipeQueryResponse{}
	mi := &file_recipe_recipe_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecipeQueryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecipeQueryResponse) ProtoMessage() {}

func (x *RecipeQueryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_recipe_recipe_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecipeQueryResponse.ProtoReflect.Descriptor instead.
func (*RecipeQueryResponse) Descriptor() ([]byte, []int) {
	return file_recipe_recipe_proto_rawDescGZIP(), []int{3}
}

func (x *RecipeQueryResponse) GetRecipes() []*GetRecipeResponse {
	if x != nil {
		return x.Recipes
	}
	return nil
}

func (x *RecipeQueryResponse) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *RecipeQueryResponse) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *RecipeQueryResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_recipe_recipe_proto protoreflect.FileDescriptor

var file_recipe_recipe_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x2f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x22, 0x2f, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x49, 0x64, 0x22, 0xb6,
	0x02, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x49,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x6e, 0x67, 0x72, 0x65,
	0x64, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x69, 0x6e,
	0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x65,
	0x70, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x65, 0x70, 0x73, 0x12,
	0x29, 0x0a, 0x10, 0x6e, 0x75, 0x74, 0x72, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x69,
	0x6e, 0x66, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6e, 0x75, 0x74, 0x72, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2d, 0x0a, 0x12, 0x61, 0x6c,
	0x6c, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x64, 0x69, 0x73, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x65, 0x72,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x67, 0x79, 0x44,
	0x69, 0x73, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61,
	0x70, 0x70, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x85, 0x01, 0x0a, 0x12, 0x52, 0x65, 0x63, 0x69,
	0x70, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22,
	0x8a, 0x01, 0x0a, 0x13, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x07, 0x72, 0x65, 0x63, 0x69, 0x70,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x72, 0x65, 0x63, 0x69, 0x70,
	0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x52, 0x07, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x32, 0x99, 0x01, 0x0a,
	0x0d, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40,
	0x0a, 0x09, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x12, 0x18, 0x2e, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x2e, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x46, 0x0a, 0x0b, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x12,
	0x1a, 0x2e, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x2e, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x65, 0x2e, 0x52, 0x65, 0x63, 0x69, 0x70, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x67, 0x65, 0x7a, 0x61, 0x2f, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x65, 0x2d, 0x62, 0x6f, 0x6f, 0x6b, 0x2d, 0x61, 0x70, 0x69, 0x2d, 0x76, 0x32,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x63,
	0x69, 0x70, 0x65, 0x3b, 0x72, 0x65, 0x63, 0x69, 0x70, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
})

var (
	file_recipe_recipe_proto_rawDescOnce sync.Once
	file_recipe_recipe_proto_rawDescData []byte
)

func file_recipe_recipe_proto_rawDescGZIP() []byte {
	file_recipe_recipe_proto_rawDescOnce.Do(func() {
		file_recipe_recipe_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_recipe_recipe_proto_rawDesc), len(file_recipe_recipe_proto_rawDesc)))
	})
	return file_recipe_recipe_proto_rawDescData
}

var file_recipe_recipe_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_recipe_recipe_proto_goTypes = []any{
	(*GetRecipeRequest)(nil),    // 0: recipe.GetRecipeRequest
	(*GetRecipeResponse)(nil),   // 1: recipe.GetRecipeResponse
	(*RecipeQueryRequest)(nil),  // 2: recipe.RecipeQueryRequest
	(*RecipeQueryResponse)(nil), // 3: recipe.RecipeQueryResponse
}
var file_recipe_recipe_proto_depIdxs = []int32{
	1, // 0: recipe.RecipeQueryResponse.recipes:type_name -> recipe.GetRecipeResponse
	0, // 1: recipe.RecipeService.GetRecipe:input_type -> recipe.GetRecipeRequest
	2, // 2: recipe.RecipeService.QueryRecipe:input_type -> recipe.RecipeQueryRequest
	1, // 3: recipe.RecipeService.GetRecipe:output_type -> recipe.GetRecipeResponse
	3, // 4: recipe.RecipeService.QueryRecipe:output_type -> recipe.RecipeQueryResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_recipe_recipe_proto_init() }
func file_recipe_recipe_proto_init() {
	if File_recipe_recipe_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_recipe_recipe_proto_rawDesc), len(file_recipe_recipe_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_recipe_recipe_proto_goTypes,
		DependencyIndexes: file_recipe_recipe_proto_depIdxs,
		MessageInfos:      file_recipe_recipe_proto_msgTypes,
	}.Build()
	File_recipe_recipe_proto = out.File
	file_recipe_recipe_proto_goTypes = nil
	file_recipe_recipe_proto_depIdxs = nil
}

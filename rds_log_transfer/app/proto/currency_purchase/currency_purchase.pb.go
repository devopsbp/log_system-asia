// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: currency_purchase.proto

package currency_purchase

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CurrencyPurchaseMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId                  *string                `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	CharacterId             *string                `protobuf:"bytes,2,opt,name=character_id,json=characterId" json:"character_id,omitempty"`
	AgencyId                *int32                 `protobuf:"zigzag32,3,opt,name=agency_id,json=agencyId" json:"agency_id,omitempty"`
	PaymentType             *int32                 `protobuf:"zigzag32,4,opt,name=payment_type,json=paymentType" json:"payment_type,omitempty"`
	OrderId                 *string                `protobuf:"bytes,5,opt,name=order_id,json=orderId" json:"order_id,omitempty"`
	ItemId                  *int32                 `protobuf:"zigzag32,6,opt,name=item_id,json=itemId" json:"item_id,omitempty"`
	ItemName                *string                `protobuf:"bytes,7,opt,name=item_name,json=itemName" json:"item_name,omitempty"`
	Count                   *int32                 `protobuf:"zigzag32,8,opt,name=count" json:"count,omitempty"`
	Price                   *int32                 `protobuf:"zigzag32,9,opt,name=price" json:"price,omitempty"`
	Step                    *int32                 `protobuf:"zigzag32,10,opt,name=step" json:"step,omitempty"`
	Url                     *string                `protobuf:"bytes,11,opt,name=url" json:"url,omitempty"`
	RequestDate             *string                `protobuf:"bytes,12,opt,name=request_date,json=requestDate" json:"request_date,omitempty"`
	RequestParams           *string                `protobuf:"bytes,13,opt,name=request_params,json=requestParams" json:"request_params,omitempty"`
	ResultCode              *int32                 `protobuf:"zigzag32,14,opt,name=result_code,json=resultCode" json:"result_code,omitempty"`
	AgencyResultCode        *string                `protobuf:"bytes,15,opt,name=agency_result_code,json=agencyResultCode" json:"agency_result_code,omitempty"`
	AgencyResultMessage     *string                `protobuf:"bytes,16,opt,name=agency_result_message,json=agencyResultMessage" json:"agency_result_message,omitempty"`
	ResponseDate            *string                `protobuf:"bytes,17,opt,name=response_date,json=responseDate" json:"response_date,omitempty"`
	StatusCode              *int32                 `protobuf:"zigzag32,18,opt,name=status_code,json=statusCode" json:"status_code,omitempty"`
	Result                  *string                `protobuf:"bytes,19,opt,name=result" json:"result,omitempty"`
	UnitPrice               *float32               `protobuf:"fixed32,20,opt,name=unit_price,json=unitPrice" json:"unit_price,omitempty"`
	Currency                *string                `protobuf:"bytes,21,opt,name=currency" json:"currency,omitempty"`
	FirstPurchaseForPayment *bool                  `protobuf:"varint,22,opt,name=first_purchase_for_payment,json=firstPurchaseForPayment" json:"first_purchase_for_payment,omitempty"`
	RequestHeaderLocation   *string                `protobuf:"bytes,23,opt,name=request_header_location,json=requestHeaderLocation" json:"request_header_location,omitempty"`
	LogTime                 *timestamppb.Timestamp `protobuf:"bytes,24,opt,name=log_time,json=logTime" json:"log_time,omitempty"`
	Uuid                    *string                `protobuf:"bytes,25,opt,name=uuid" json:"uuid,omitempty"`
}

func (x *CurrencyPurchaseMsg) Reset() {
	*x = CurrencyPurchaseMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_currency_purchase_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CurrencyPurchaseMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CurrencyPurchaseMsg) ProtoMessage() {}

func (x *CurrencyPurchaseMsg) ProtoReflect() protoreflect.Message {
	mi := &file_currency_purchase_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CurrencyPurchaseMsg.ProtoReflect.Descriptor instead.
func (*CurrencyPurchaseMsg) Descriptor() ([]byte, []int) {
	return file_currency_purchase_proto_rawDescGZIP(), []int{0}
}

func (x *CurrencyPurchaseMsg) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetCharacterId() string {
	if x != nil && x.CharacterId != nil {
		return *x.CharacterId
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetAgencyId() int32 {
	if x != nil && x.AgencyId != nil {
		return *x.AgencyId
	}
	return 0
}

func (x *CurrencyPurchaseMsg) GetPaymentType() int32 {
	if x != nil && x.PaymentType != nil {
		return *x.PaymentType
	}
	return 0
}

func (x *CurrencyPurchaseMsg) GetOrderId() string {
	if x != nil && x.OrderId != nil {
		return *x.OrderId
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetItemId() int32 {
	if x != nil && x.ItemId != nil {
		return *x.ItemId
	}
	return 0
}

func (x *CurrencyPurchaseMsg) GetItemName() string {
	if x != nil && x.ItemName != nil {
		return *x.ItemName
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetCount() int32 {
	if x != nil && x.Count != nil {
		return *x.Count
	}
	return 0
}

func (x *CurrencyPurchaseMsg) GetPrice() int32 {
	if x != nil && x.Price != nil {
		return *x.Price
	}
	return 0
}

func (x *CurrencyPurchaseMsg) GetStep() int32 {
	if x != nil && x.Step != nil {
		return *x.Step
	}
	return 0
}

func (x *CurrencyPurchaseMsg) GetUrl() string {
	if x != nil && x.Url != nil {
		return *x.Url
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetRequestDate() string {
	if x != nil && x.RequestDate != nil {
		return *x.RequestDate
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetRequestParams() string {
	if x != nil && x.RequestParams != nil {
		return *x.RequestParams
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetResultCode() int32 {
	if x != nil && x.ResultCode != nil {
		return *x.ResultCode
	}
	return 0
}

func (x *CurrencyPurchaseMsg) GetAgencyResultCode() string {
	if x != nil && x.AgencyResultCode != nil {
		return *x.AgencyResultCode
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetAgencyResultMessage() string {
	if x != nil && x.AgencyResultMessage != nil {
		return *x.AgencyResultMessage
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetResponseDate() string {
	if x != nil && x.ResponseDate != nil {
		return *x.ResponseDate
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetStatusCode() int32 {
	if x != nil && x.StatusCode != nil {
		return *x.StatusCode
	}
	return 0
}

func (x *CurrencyPurchaseMsg) GetResult() string {
	if x != nil && x.Result != nil {
		return *x.Result
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetUnitPrice() float32 {
	if x != nil && x.UnitPrice != nil {
		return *x.UnitPrice
	}
	return 0
}

func (x *CurrencyPurchaseMsg) GetCurrency() string {
	if x != nil && x.Currency != nil {
		return *x.Currency
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetFirstPurchaseForPayment() bool {
	if x != nil && x.FirstPurchaseForPayment != nil {
		return *x.FirstPurchaseForPayment
	}
	return false
}

func (x *CurrencyPurchaseMsg) GetRequestHeaderLocation() string {
	if x != nil && x.RequestHeaderLocation != nil {
		return *x.RequestHeaderLocation
	}
	return ""
}

func (x *CurrencyPurchaseMsg) GetLogTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LogTime
	}
	return nil
}

func (x *CurrencyPurchaseMsg) GetUuid() string {
	if x != nil && x.Uuid != nil {
		return *x.Uuid
	}
	return ""
}

var File_currency_purchase_proto protoreflect.FileDescriptor

var file_currency_purchase_proto_rawDesc = []byte{
	0x0a, 0x17, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x70, 0x75, 0x72, 0x63, 0x68,
	0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x63, 0x79, 0x5f, 0x70, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xda, 0x06,
	0x0a, 0x13, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61,
	0x73, 0x65, 0x4d, 0x73, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x21,
	0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x67, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x11, 0x52, 0x08, 0x61, 0x67, 0x65, 0x6e, 0x63, 0x79, 0x49, 0x64, 0x12, 0x21,
	0x0a, 0x0c, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x11, 0x52, 0x0b, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x11, 0x52, 0x06, 0x69,
	0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x11, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x11, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x74, 0x65, 0x70, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x11, 0x52, 0x04, 0x73, 0x74,
	0x65, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x11, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x2c, 0x0a, 0x12, 0x61, 0x67, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x61, 0x67, 0x65,
	0x6e, 0x63, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x32, 0x0a,
	0x15, 0x61, 0x67, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x61, 0x67,
	0x65, 0x6e, 0x63, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x64, 0x61,
	0x74, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x11, 0x52, 0x0a, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x14, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x09, 0x75, 0x6e, 0x69, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x3b, 0x0a, 0x1a, 0x66, 0x69,
	0x72, 0x73, 0x74, 0x5f, 0x70, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x5f, 0x66, 0x6f, 0x72,
	0x5f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x16, 0x20, 0x01, 0x28, 0x08, 0x52, 0x17,
	0x66, 0x69, 0x72, 0x73, 0x74, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x46, 0x6f, 0x72,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x36, 0x0a, 0x17, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x17, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x35, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x18, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x6c,
	0x6f, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x19,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x42, 0x19, 0x5a, 0x17, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x70, 0x75, 0x72,
	0x63, 0x68, 0x61, 0x73, 0x65,
}

var (
	file_currency_purchase_proto_rawDescOnce sync.Once
	file_currency_purchase_proto_rawDescData = file_currency_purchase_proto_rawDesc
)

func file_currency_purchase_proto_rawDescGZIP() []byte {
	file_currency_purchase_proto_rawDescOnce.Do(func() {
		file_currency_purchase_proto_rawDescData = protoimpl.X.CompressGZIP(file_currency_purchase_proto_rawDescData)
	})
	return file_currency_purchase_proto_rawDescData
}

var file_currency_purchase_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_currency_purchase_proto_goTypes = []interface{}{
	(*CurrencyPurchaseMsg)(nil),   // 0: currency_purchase.CurrencyPurchaseMsg
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_currency_purchase_proto_depIdxs = []int32{
	1, // 0: currency_purchase.CurrencyPurchaseMsg.log_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_currency_purchase_proto_init() }
func file_currency_purchase_proto_init() {
	if File_currency_purchase_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_currency_purchase_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CurrencyPurchaseMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_currency_purchase_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_currency_purchase_proto_goTypes,
		DependencyIndexes: file_currency_purchase_proto_depIdxs,
		MessageInfos:      file_currency_purchase_proto_msgTypes,
	}.Build()
	File_currency_purchase_proto = out.File
	file_currency_purchase_proto_rawDesc = nil
	file_currency_purchase_proto_goTypes = nil
	file_currency_purchase_proto_depIdxs = nil
}
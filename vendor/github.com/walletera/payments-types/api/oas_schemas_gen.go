// Code generated by ogen, DO NOT EDIT.

package api

import (
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

// Bank account details.
// Ref: #/components/schemas/accountDetails
type AccountDetails struct {
	InstitutionName OptString                    `json:"institutionName"`
	InstitutionId   OptString                    `json:"institutionId"`
	Currency        OptString                    `json:"currency"`
	AccountType     OptAccountDetailsAccountType `json:"accountType"`
	// Extra account details. The details depend on the accountType.
	AccountDetails OptAccountDetailsAccountDetails `json:"accountDetails"`
}

// GetInstitutionName returns the value of InstitutionName.
func (s *AccountDetails) GetInstitutionName() OptString {
	return s.InstitutionName
}

// GetInstitutionId returns the value of InstitutionId.
func (s *AccountDetails) GetInstitutionId() OptString {
	return s.InstitutionId
}

// GetCurrency returns the value of Currency.
func (s *AccountDetails) GetCurrency() OptString {
	return s.Currency
}

// GetAccountType returns the value of AccountType.
func (s *AccountDetails) GetAccountType() OptAccountDetailsAccountType {
	return s.AccountType
}

// GetAccountDetails returns the value of AccountDetails.
func (s *AccountDetails) GetAccountDetails() OptAccountDetailsAccountDetails {
	return s.AccountDetails
}

// SetInstitutionName sets the value of InstitutionName.
func (s *AccountDetails) SetInstitutionName(val OptString) {
	s.InstitutionName = val
}

// SetInstitutionId sets the value of InstitutionId.
func (s *AccountDetails) SetInstitutionId(val OptString) {
	s.InstitutionId = val
}

// SetCurrency sets the value of Currency.
func (s *AccountDetails) SetCurrency(val OptString) {
	s.Currency = val
}

// SetAccountType sets the value of AccountType.
func (s *AccountDetails) SetAccountType(val OptAccountDetailsAccountType) {
	s.AccountType = val
}

// SetAccountDetails sets the value of AccountDetails.
func (s *AccountDetails) SetAccountDetails(val OptAccountDetailsAccountDetails) {
	s.AccountDetails = val
}

// Extra account details. The details depend on the accountType.
type AccountDetailsAccountDetails struct {
	OneOf AccountDetailsAccountDetailsSum
}

// GetOneOf returns the value of OneOf.
func (s *AccountDetailsAccountDetails) GetOneOf() AccountDetailsAccountDetailsSum {
	return s.OneOf
}

// SetOneOf sets the value of OneOf.
func (s *AccountDetailsAccountDetails) SetOneOf(val AccountDetailsAccountDetailsSum) {
	s.OneOf = val
}

// AccountDetailsAccountDetailsSum represents sum type.
type AccountDetailsAccountDetailsSum struct {
	Type                  AccountDetailsAccountDetailsSumType // switch on this field
	CvuAccountDetails     CvuAccountDetails
	DinopayAccountDetails DinopayAccountDetails
}

// AccountDetailsAccountDetailsSumType is oneOf type of AccountDetailsAccountDetailsSum.
type AccountDetailsAccountDetailsSumType string

// Possible values for AccountDetailsAccountDetailsSumType.
const (
	CvuAccountDetailsAccountDetailsAccountDetailsSum     AccountDetailsAccountDetailsSumType = "CvuAccountDetails"
	DinopayAccountDetailsAccountDetailsAccountDetailsSum AccountDetailsAccountDetailsSumType = "DinopayAccountDetails"
)

// IsCvuAccountDetails reports whether AccountDetailsAccountDetailsSum is CvuAccountDetails.
func (s AccountDetailsAccountDetailsSum) IsCvuAccountDetails() bool {
	return s.Type == CvuAccountDetailsAccountDetailsAccountDetailsSum
}

// IsDinopayAccountDetails reports whether AccountDetailsAccountDetailsSum is DinopayAccountDetails.
func (s AccountDetailsAccountDetailsSum) IsDinopayAccountDetails() bool {
	return s.Type == DinopayAccountDetailsAccountDetailsAccountDetailsSum
}

// SetCvuAccountDetails sets AccountDetailsAccountDetailsSum to CvuAccountDetails.
func (s *AccountDetailsAccountDetailsSum) SetCvuAccountDetails(v CvuAccountDetails) {
	s.Type = CvuAccountDetailsAccountDetailsAccountDetailsSum
	s.CvuAccountDetails = v
}

// GetCvuAccountDetails returns CvuAccountDetails and true boolean if AccountDetailsAccountDetailsSum is CvuAccountDetails.
func (s AccountDetailsAccountDetailsSum) GetCvuAccountDetails() (v CvuAccountDetails, ok bool) {
	if !s.IsCvuAccountDetails() {
		return v, false
	}
	return s.CvuAccountDetails, true
}

// NewCvuAccountDetailsAccountDetailsAccountDetailsSum returns new AccountDetailsAccountDetailsSum from CvuAccountDetails.
func NewCvuAccountDetailsAccountDetailsAccountDetailsSum(v CvuAccountDetails) AccountDetailsAccountDetailsSum {
	var s AccountDetailsAccountDetailsSum
	s.SetCvuAccountDetails(v)
	return s
}

// SetDinopayAccountDetails sets AccountDetailsAccountDetailsSum to DinopayAccountDetails.
func (s *AccountDetailsAccountDetailsSum) SetDinopayAccountDetails(v DinopayAccountDetails) {
	s.Type = DinopayAccountDetailsAccountDetailsAccountDetailsSum
	s.DinopayAccountDetails = v
}

// GetDinopayAccountDetails returns DinopayAccountDetails and true boolean if AccountDetailsAccountDetailsSum is DinopayAccountDetails.
func (s AccountDetailsAccountDetailsSum) GetDinopayAccountDetails() (v DinopayAccountDetails, ok bool) {
	if !s.IsDinopayAccountDetails() {
		return v, false
	}
	return s.DinopayAccountDetails, true
}

// NewDinopayAccountDetailsAccountDetailsAccountDetailsSum returns new AccountDetailsAccountDetailsSum from DinopayAccountDetails.
func NewDinopayAccountDetailsAccountDetailsAccountDetailsSum(v DinopayAccountDetails) AccountDetailsAccountDetailsSum {
	var s AccountDetailsAccountDetailsSum
	s.SetDinopayAccountDetails(v)
	return s
}

type AccountDetailsAccountType string

const (
	AccountDetailsAccountTypeCvu     AccountDetailsAccountType = "cvu"
	AccountDetailsAccountTypeDinopay AccountDetailsAccountType = "dinopay"
)

// AllValues returns all AccountDetailsAccountType values.
func (AccountDetailsAccountType) AllValues() []AccountDetailsAccountType {
	return []AccountDetailsAccountType{
		AccountDetailsAccountTypeCvu,
		AccountDetailsAccountTypeDinopay,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s AccountDetailsAccountType) MarshalText() ([]byte, error) {
	switch s {
	case AccountDetailsAccountTypeCvu:
		return []byte(s), nil
	case AccountDetailsAccountTypeDinopay:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *AccountDetailsAccountType) UnmarshalText(data []byte) error {
	switch AccountDetailsAccountType(data) {
	case AccountDetailsAccountTypeCvu:
		*s = AccountDetailsAccountTypeCvu
		return nil
	case AccountDetailsAccountTypeDinopay:
		*s = AccountDetailsAccountTypeDinopay
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Body of the error responses.
// Ref: #/components/schemas/apiError
type ApiError struct {
	// A message describing the error.
	ErrorMessage string `json:"errorMessage"`
	// A unique identifier for the specific error.
	ErrorCode uuid.UUID `json:"errorCode"`
}

// GetErrorMessage returns the value of ErrorMessage.
func (s *ApiError) GetErrorMessage() string {
	return s.ErrorMessage
}

// GetErrorCode returns the value of ErrorCode.
func (s *ApiError) GetErrorCode() uuid.UUID {
	return s.ErrorCode
}

// SetErrorMessage sets the value of ErrorMessage.
func (s *ApiError) SetErrorMessage(val string) {
	s.ErrorMessage = val
}

// SetErrorCode sets the value of ErrorCode.
func (s *ApiError) SetErrorCode(val uuid.UUID) {
	s.ErrorCode = val
}

type BearerAuth struct {
	Token string
}

// GetToken returns the value of Token.
func (s *BearerAuth) GetToken() string {
	return s.Token
}

// SetToken sets the value of Token.
func (s *BearerAuth) SetToken(val string) {
	s.Token = val
}

// Ref: #/components/schemas/cvuAccountDetails
type CvuAccountDetails struct {
	// Account owner national identification number.
	Cuit  OptString `json:"cuit"`
	Cvu   OptString `json:"cvu"`
	Alias OptString `json:"alias"`
}

// GetCuit returns the value of Cuit.
func (s *CvuAccountDetails) GetCuit() OptString {
	return s.Cuit
}

// GetCvu returns the value of Cvu.
func (s *CvuAccountDetails) GetCvu() OptString {
	return s.Cvu
}

// GetAlias returns the value of Alias.
func (s *CvuAccountDetails) GetAlias() OptString {
	return s.Alias
}

// SetCuit sets the value of Cuit.
func (s *CvuAccountDetails) SetCuit(val OptString) {
	s.Cuit = val
}

// SetCvu sets the value of Cvu.
func (s *CvuAccountDetails) SetCvu(val OptString) {
	s.Cvu = val
}

// SetAlias sets the value of Alias.
func (s *CvuAccountDetails) SetAlias(val OptString) {
	s.Alias = val
}

// Ref: #/components/schemas/dinopayAccountDetails
type DinopayAccountDetails struct {
	// Name of the owner of the account.
	AccountHolder string `json:"accountHolder"`
	// Account number on DinoPay.
	AccountNumber string `json:"accountNumber"`
}

// GetAccountHolder returns the value of AccountHolder.
func (s *DinopayAccountDetails) GetAccountHolder() string {
	return s.AccountHolder
}

// GetAccountNumber returns the value of AccountNumber.
func (s *DinopayAccountDetails) GetAccountNumber() string {
	return s.AccountNumber
}

// SetAccountHolder sets the value of AccountHolder.
func (s *DinopayAccountDetails) SetAccountHolder(val string) {
	s.AccountHolder = val
}

// SetAccountNumber sets the value of AccountNumber.
func (s *DinopayAccountDetails) SetAccountNumber(val string) {
	s.AccountNumber = val
}

// GetPaymentInternalServerError is response for GetPayment operation.
type GetPaymentInternalServerError struct{}

func (*GetPaymentInternalServerError) getPaymentRes() {}

// GetPaymentNotFound is response for GetPayment operation.
type GetPaymentNotFound struct{}

func (*GetPaymentNotFound) getPaymentRes() {}

// GetPaymentUnauthorized is response for GetPayment operation.
type GetPaymentUnauthorized struct{}

func (*GetPaymentUnauthorized) getPaymentRes() {}

// NewOptAccountDetails returns new OptAccountDetails with value set to v.
func NewOptAccountDetails(v AccountDetails) OptAccountDetails {
	return OptAccountDetails{
		Value: v,
		Set:   true,
	}
}

// OptAccountDetails is optional AccountDetails.
type OptAccountDetails struct {
	Value AccountDetails
	Set   bool
}

// IsSet returns true if OptAccountDetails was set.
func (o OptAccountDetails) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptAccountDetails) Reset() {
	var v AccountDetails
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptAccountDetails) SetTo(v AccountDetails) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptAccountDetails) Get() (v AccountDetails, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptAccountDetails) Or(d AccountDetails) AccountDetails {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptAccountDetailsAccountDetails returns new OptAccountDetailsAccountDetails with value set to v.
func NewOptAccountDetailsAccountDetails(v AccountDetailsAccountDetails) OptAccountDetailsAccountDetails {
	return OptAccountDetailsAccountDetails{
		Value: v,
		Set:   true,
	}
}

// OptAccountDetailsAccountDetails is optional AccountDetailsAccountDetails.
type OptAccountDetailsAccountDetails struct {
	Value AccountDetailsAccountDetails
	Set   bool
}

// IsSet returns true if OptAccountDetailsAccountDetails was set.
func (o OptAccountDetailsAccountDetails) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptAccountDetailsAccountDetails) Reset() {
	var v AccountDetailsAccountDetails
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptAccountDetailsAccountDetails) SetTo(v AccountDetailsAccountDetails) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptAccountDetailsAccountDetails) Get() (v AccountDetailsAccountDetails, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptAccountDetailsAccountDetails) Or(d AccountDetailsAccountDetails) AccountDetailsAccountDetails {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptAccountDetailsAccountType returns new OptAccountDetailsAccountType with value set to v.
func NewOptAccountDetailsAccountType(v AccountDetailsAccountType) OptAccountDetailsAccountType {
	return OptAccountDetailsAccountType{
		Value: v,
		Set:   true,
	}
}

// OptAccountDetailsAccountType is optional AccountDetailsAccountType.
type OptAccountDetailsAccountType struct {
	Value AccountDetailsAccountType
	Set   bool
}

// IsSet returns true if OptAccountDetailsAccountType was set.
func (o OptAccountDetailsAccountType) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptAccountDetailsAccountType) Reset() {
	var v AccountDetailsAccountType
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptAccountDetailsAccountType) SetTo(v AccountDetailsAccountType) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptAccountDetailsAccountType) Get() (v AccountDetailsAccountType, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptAccountDetailsAccountType) Or(d AccountDetailsAccountType) AccountDetailsAccountType {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptDateTime returns new OptDateTime with value set to v.
func NewOptDateTime(v time.Time) OptDateTime {
	return OptDateTime{
		Value: v,
		Set:   true,
	}
}

// OptDateTime is optional time.Time.
type OptDateTime struct {
	Value time.Time
	Set   bool
}

// IsSet returns true if OptDateTime was set.
func (o OptDateTime) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptDateTime) Reset() {
	var v time.Time
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptDateTime) SetTo(v time.Time) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptDateTime) Get() (v time.Time, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptDateTime) Or(d time.Time) time.Time {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptPaymentDirection returns new OptPaymentDirection with value set to v.
func NewOptPaymentDirection(v PaymentDirection) OptPaymentDirection {
	return OptPaymentDirection{
		Value: v,
		Set:   true,
	}
}

// OptPaymentDirection is optional PaymentDirection.
type OptPaymentDirection struct {
	Value PaymentDirection
	Set   bool
}

// IsSet returns true if OptPaymentDirection was set.
func (o OptPaymentDirection) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptPaymentDirection) Reset() {
	var v PaymentDirection
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptPaymentDirection) SetTo(v PaymentDirection) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptPaymentDirection) Get() (v PaymentDirection, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptPaymentDirection) Or(d PaymentDirection) PaymentDirection {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptPaymentStatus returns new OptPaymentStatus with value set to v.
func NewOptPaymentStatus(v PaymentStatus) OptPaymentStatus {
	return OptPaymentStatus{
		Value: v,
		Set:   true,
	}
}

// OptPaymentStatus is optional PaymentStatus.
type OptPaymentStatus struct {
	Value PaymentStatus
	Set   bool
}

// IsSet returns true if OptPaymentStatus was set.
func (o OptPaymentStatus) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptPaymentStatus) Reset() {
	var v PaymentStatus
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptPaymentStatus) SetTo(v PaymentStatus) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptPaymentStatus) Get() (v PaymentStatus, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptPaymentStatus) Or(d PaymentStatus) PaymentStatus {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptUUID returns new OptUUID with value set to v.
func NewOptUUID(v uuid.UUID) OptUUID {
	return OptUUID{
		Value: v,
		Set:   true,
	}
}

// OptUUID is optional uuid.UUID.
type OptUUID struct {
	Value uuid.UUID
	Set   bool
}

// IsSet returns true if OptUUID was set.
func (o OptUUID) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptUUID) Reset() {
	var v uuid.UUID
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptUUID) SetTo(v uuid.UUID) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptUUID) Get() (v uuid.UUID, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptUUID) Or(d uuid.UUID) uuid.UUID {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

type PatchPaymentBadRequest ApiError

func (*PatchPaymentBadRequest) patchPaymentRes() {}

type PatchPaymentInternalServerError ApiError

func (*PatchPaymentInternalServerError) patchPaymentRes() {}

// PatchPaymentOK is response for PatchPayment operation.
type PatchPaymentOK struct{}

func (*PatchPaymentOK) patchPaymentRes() {}

type PatchPaymentUnauthorized ApiError

func (*PatchPaymentUnauthorized) patchPaymentRes() {}

// Payment type.
// Ref: #/components/schemas/payment
type Payment struct {
	// Payment id.
	ID uuid.UUID `json:"id"`
	// Payment amount.
	Amount float64 `json:"amount"`
	// Payment currency.
	Currency  string              `json:"currency"`
	Direction OptPaymentDirection `json:"direction"`
	// The customer that sent or received the payment.
	CustomerId OptUUID `json:"customerId"`
	// Id assigned to the operation by an external payment provider.
	ExternalId OptString `json:"externalId"`
	// Id assigned to the operation by the payment scheme or clearing institution.
	SchemeId    OptString         `json:"schemeId"`
	Beneficiary OptAccountDetails `json:"beneficiary"`
	Debtor      OptAccountDetails `json:"debtor"`
	Status      OptPaymentStatus  `json:"status"`
	CreatedAt   OptDateTime       `json:"createdAt"`
	UpdatedAt   OptDateTime       `json:"updatedAt"`
}

// GetID returns the value of ID.
func (s *Payment) GetID() uuid.UUID {
	return s.ID
}

// GetAmount returns the value of Amount.
func (s *Payment) GetAmount() float64 {
	return s.Amount
}

// GetCurrency returns the value of Currency.
func (s *Payment) GetCurrency() string {
	return s.Currency
}

// GetDirection returns the value of Direction.
func (s *Payment) GetDirection() OptPaymentDirection {
	return s.Direction
}

// GetCustomerId returns the value of CustomerId.
func (s *Payment) GetCustomerId() OptUUID {
	return s.CustomerId
}

// GetExternalId returns the value of ExternalId.
func (s *Payment) GetExternalId() OptString {
	return s.ExternalId
}

// GetSchemeId returns the value of SchemeId.
func (s *Payment) GetSchemeId() OptString {
	return s.SchemeId
}

// GetBeneficiary returns the value of Beneficiary.
func (s *Payment) GetBeneficiary() OptAccountDetails {
	return s.Beneficiary
}

// GetDebtor returns the value of Debtor.
func (s *Payment) GetDebtor() OptAccountDetails {
	return s.Debtor
}

// GetStatus returns the value of Status.
func (s *Payment) GetStatus() OptPaymentStatus {
	return s.Status
}

// GetCreatedAt returns the value of CreatedAt.
func (s *Payment) GetCreatedAt() OptDateTime {
	return s.CreatedAt
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *Payment) GetUpdatedAt() OptDateTime {
	return s.UpdatedAt
}

// SetID sets the value of ID.
func (s *Payment) SetID(val uuid.UUID) {
	s.ID = val
}

// SetAmount sets the value of Amount.
func (s *Payment) SetAmount(val float64) {
	s.Amount = val
}

// SetCurrency sets the value of Currency.
func (s *Payment) SetCurrency(val string) {
	s.Currency = val
}

// SetDirection sets the value of Direction.
func (s *Payment) SetDirection(val OptPaymentDirection) {
	s.Direction = val
}

// SetCustomerId sets the value of CustomerId.
func (s *Payment) SetCustomerId(val OptUUID) {
	s.CustomerId = val
}

// SetExternalId sets the value of ExternalId.
func (s *Payment) SetExternalId(val OptString) {
	s.ExternalId = val
}

// SetSchemeId sets the value of SchemeId.
func (s *Payment) SetSchemeId(val OptString) {
	s.SchemeId = val
}

// SetBeneficiary sets the value of Beneficiary.
func (s *Payment) SetBeneficiary(val OptAccountDetails) {
	s.Beneficiary = val
}

// SetDebtor sets the value of Debtor.
func (s *Payment) SetDebtor(val OptAccountDetails) {
	s.Debtor = val
}

// SetStatus sets the value of Status.
func (s *Payment) SetStatus(val OptPaymentStatus) {
	s.Status = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *Payment) SetCreatedAt(val OptDateTime) {
	s.CreatedAt = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *Payment) SetUpdatedAt(val OptDateTime) {
	s.UpdatedAt = val
}

func (*Payment) getPaymentRes()  {}
func (*Payment) postPaymentRes() {}

type PaymentDirection string

const (
	PaymentDirectionInbound  PaymentDirection = "inbound"
	PaymentDirectionOutbound PaymentDirection = "outbound"
)

// AllValues returns all PaymentDirection values.
func (PaymentDirection) AllValues() []PaymentDirection {
	return []PaymentDirection{
		PaymentDirectionInbound,
		PaymentDirectionOutbound,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s PaymentDirection) MarshalText() ([]byte, error) {
	switch s {
	case PaymentDirectionInbound:
		return []byte(s), nil
	case PaymentDirectionOutbound:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *PaymentDirection) UnmarshalText(data []byte) error {
	switch PaymentDirection(data) {
	case PaymentDirectionInbound:
		*s = PaymentDirectionInbound
		return nil
	case PaymentDirectionOutbound:
		*s = PaymentDirectionOutbound
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Ref: #/components/schemas/paymentStatus
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusDelivered PaymentStatus = "delivered"
	PaymentStatusConfirmed PaymentStatus = "confirmed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRejected  PaymentStatus = "rejected"
)

// AllValues returns all PaymentStatus values.
func (PaymentStatus) AllValues() []PaymentStatus {
	return []PaymentStatus{
		PaymentStatusPending,
		PaymentStatusDelivered,
		PaymentStatusConfirmed,
		PaymentStatusFailed,
		PaymentStatusRejected,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s PaymentStatus) MarshalText() ([]byte, error) {
	switch s {
	case PaymentStatusPending:
		return []byte(s), nil
	case PaymentStatusDelivered:
		return []byte(s), nil
	case PaymentStatusConfirmed:
		return []byte(s), nil
	case PaymentStatusFailed:
		return []byte(s), nil
	case PaymentStatusRejected:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *PaymentStatus) UnmarshalText(data []byte) error {
	switch PaymentStatus(data) {
	case PaymentStatusPending:
		*s = PaymentStatusPending
		return nil
	case PaymentStatusDelivered:
		*s = PaymentStatusDelivered
		return nil
	case PaymentStatusConfirmed:
		*s = PaymentStatusConfirmed
		return nil
	case PaymentStatusFailed:
		*s = PaymentStatusFailed
		return nil
	case PaymentStatusRejected:
		*s = PaymentStatusRejected
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Body of the PATH /withdrawal request.
// Ref: #/components/schemas/paymentUpdate
type PaymentUpdate struct {
	// Payment Id.
	PaymentId uuid.UUID `json:"paymentId"`
	// Id assigned to the operation by the external payment provider.
	ExternalId OptString `json:"externalId"`
	// Payment status.
	Status PaymentStatus `json:"status"`
}

// GetPaymentId returns the value of PaymentId.
func (s *PaymentUpdate) GetPaymentId() uuid.UUID {
	return s.PaymentId
}

// GetExternalId returns the value of ExternalId.
func (s *PaymentUpdate) GetExternalId() OptString {
	return s.ExternalId
}

// GetStatus returns the value of Status.
func (s *PaymentUpdate) GetStatus() PaymentStatus {
	return s.Status
}

// SetPaymentId sets the value of PaymentId.
func (s *PaymentUpdate) SetPaymentId(val uuid.UUID) {
	s.PaymentId = val
}

// SetExternalId sets the value of ExternalId.
func (s *PaymentUpdate) SetExternalId(val OptString) {
	s.ExternalId = val
}

// SetStatus sets the value of Status.
func (s *PaymentUpdate) SetStatus(val PaymentStatus) {
	s.Status = val
}

type PostPaymentBadRequest ApiError

func (*PostPaymentBadRequest) postPaymentRes() {}

type PostPaymentConflict ApiError

func (*PostPaymentConflict) postPaymentRes() {}

type PostPaymentInternalServerError ApiError

func (*PostPaymentInternalServerError) postPaymentRes() {}

type PostPaymentUnauthorized ApiError

func (*PostPaymentUnauthorized) postPaymentRes() {}

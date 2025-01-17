// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpzap

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/wrp-go/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestObserver_ObserveWRP(t *testing.T) {
	num := int64(123)
	headers := []string{"header1", "header2"}
	payload := []byte("test payload")

	tests := []struct {
		name            string
		nilLogger       bool
		fields          []FieldOpt
		input_message   wrp.Message
		expected_fields []zap.Field
	}{
		{
			name:      "nil logger",
			nilLogger: true,
		}, {
			name:            "log message type as default (num)",
			fields:          []FieldOpt{LogMessageType()},
			input_message:   wrp.Message{Type: wrp.SimpleRequestResponseMessageType},
			expected_fields: []zap.Field{zap.Int(fMsgType, int(wrp.SimpleRequestResponseMessageType))},
		}, {
			name:            "log message type as num",
			fields:          []FieldOpt{LogMessageTypeAsNum()},
			input_message:   wrp.Message{Type: wrp.SimpleRequestResponseMessageType},
			expected_fields: []zap.Field{zap.Int(fMsgType, int(wrp.SimpleRequestResponseMessageType))},
		}, {
			name:            "log message type as string",
			fields:          []FieldOpt{LogMessageTypeAsString()},
			input_message:   wrp.Message{Type: wrp.SimpleRequestResponseMessageType},
			expected_fields: []zap.Field{zap.Stringer(fMsgType, wrp.SimpleRequestResponseMessageType)},
		}, {
			name:            "log source",
			fields:          []FieldOpt{LogSource()},
			input_message:   wrp.Message{Source: "test source"},
			expected_fields: []zap.Field{zap.String(fSource, "test source")},
		}, {
			name:            "log destination",
			fields:          []FieldOpt{LogDestination()},
			input_message:   wrp.Message{Destination: "test destination"},
			expected_fields: []zap.Field{zap.String(fDestination, "test destination")},
		}, {
			name:            "log transaction uuid",
			fields:          []FieldOpt{LogTransactionUUID()},
			input_message:   wrp.Message{TransactionUUID: "test uuid"},
			expected_fields: []zap.Field{zap.String(fTransactionUUID, "test uuid")},
		}, {
			name:            "log content type",
			fields:          []FieldOpt{LogContentType()},
			input_message:   wrp.Message{ContentType: "test content type"},
			expected_fields: []zap.Field{zap.String(fContentType, "test content type")},
		}, {
			name:            "log accept",
			fields:          []FieldOpt{LogAccept()},
			input_message:   wrp.Message{Accept: "test accept"},
			expected_fields: []zap.Field{zap.String(fAccept, "test accept")},
		}, {
			name:            "log nil status",
			fields:          []FieldOpt{LogStatus()},
			input_message:   wrp.Message{Status: nil},
			expected_fields: []zap.Field{zap.Int64p(fStatus, nil)},
		}, {
			name:            "log num status",
			fields:          []FieldOpt{LogStatus()},
			input_message:   wrp.Message{Status: &num},
			expected_fields: []zap.Field{zap.Int64p(fStatus, &num)},
		}, {
			name:            "log nil rdr",
			fields:          []FieldOpt{LogRequestDeliveryResponse()},
			input_message:   wrp.Message{RequestDeliveryResponse: nil},
			expected_fields: []zap.Field{zap.Int64p(fRequestDeliveryResponse, nil)},
		}, {
			name:            "log num rdr",
			fields:          []FieldOpt{LogRequestDeliveryResponse()},
			input_message:   wrp.Message{RequestDeliveryResponse: &num},
			expected_fields: []zap.Field{zap.Int64p(fRequestDeliveryResponse, &num)},
		}, {
			name:            "log headers",
			fields:          []FieldOpt{LogHeaders()},
			input_message:   wrp.Message{Headers: headers},
			expected_fields: []zap.Field{zap.Strings(fHeaders, headers)},
		}, {
			name:            "log metadata",
			fields:          []FieldOpt{LogMetadata()},
			input_message:   wrp.Message{Metadata: map[string]string{"key": "value"}},
			expected_fields: []zap.Field{zap.Any(fMetadata, map[string]string{"key": "value"})},
		}, {
			name:            "log path",
			fields:          []FieldOpt{LogPath()},
			input_message:   wrp.Message{Path: "test path"},
			expected_fields: []zap.Field{zap.String(fPath, "test path")},
		}, {
			name:            "log payload",
			fields:          []FieldOpt{LogPayload()},
			input_message:   wrp.Message{Payload: payload},
			expected_fields: []zap.Field{zap.Binary(fPayload, payload)},
		}, {
			name:            "log payload size",
			fields:          []FieldOpt{LogPayloadSize()},
			input_message:   wrp.Message{Payload: payload},
			expected_fields: []zap.Field{zap.Int(fPayloadSize, len(payload))},
		}, {
			name:            "log service name",
			fields:          []FieldOpt{LogServiceName()},
			input_message:   wrp.Message{ServiceName: "test service name"},
			expected_fields: []zap.Field{zap.String(fServiceName, "test service name")},
		}, {
			name:            "log url",
			fields:          []FieldOpt{LogURL()},
			input_message:   wrp.Message{URL: "test url"},
			expected_fields: []zap.Field{zap.String(fURL, "test url")},
		}, {
			name:            "log partner IDs",
			fields:          []FieldOpt{LogPartnerIDs()},
			input_message:   wrp.Message{PartnerIDs: []string{"partner1", "partner2"}},
			expected_fields: []zap.Field{zap.Strings(fPartnerIDs, []string{"partner1", "partner2"})},
		}, {
			name:            "log session ID",
			fields:          []FieldOpt{LogSessionID()},
			input_message:   wrp.Message{SessionID: "session123"},
			expected_fields: []zap.Field{zap.String(fSessionID, "session123")},
		}, {
			name:            "log quality of service",
			fields:          []FieldOpt{LogQualityOfService()},
			input_message:   wrp.Message{QualityOfService: 2},
			expected_fields: []zap.Field{zap.Int(fQualityOfService, 2)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			text := "test message"
			ob := Observer{
				Logger:  logger,
				Level:   zap.InfoLevel,
				Message: text,
				Fields:  tt.fields,
			}

			if tt.nilLogger {
				ob.Logger = nil
			}

			ob.ObserveWRP(tt.input_message)

			entries := recorded.All()
			if tt.nilLogger {
				require.Len(t, entries, 0)
				return
			}
			require.Len(t, entries, 1)
			entry := entries[0]
			assert.Equal(t, zap.InfoLevel, entry.Level)
			assert.Equal(t, text, entry.Message)
			assert.ElementsMatch(t, tt.expected_fields, entry.Context)
		})
	}
}

func TestFieldOpt_JSONTags(t *testing.T) {
	fieldMap := map[string]string{
		"Type":                    fMsgType,
		"Source":                  fSource,
		"Destination":             fDestination,
		"TransactionUUID":         fTransactionUUID,
		"ContentType":             fContentType,
		"Accept":                  fAccept,
		"Status":                  fStatus,
		"RequestDeliveryResponse": fRequestDeliveryResponse,
		"Headers":                 fHeaders,
		"Metadata":                fMetadata,
		"Path":                    fPath,
		"Payload":                 fPayload,
		"ServiceName":             fServiceName,
		"URL":                     fURL,
		"PartnerIDs":              fPartnerIDs,
		"SessionID":               fSessionID,
		"QualityOfService":        fQualityOfService,
	}

	ignored := map[string]struct{}{
		"Spans":        {},
		"IncludeSpans": {},
	}

	msgType := reflect.TypeOf(wrp.Message{})

	for fieldName, expectedTag := range fieldMap {
		t.Run(fieldName, func(t *testing.T) {
			field, found := msgType.FieldByName(fieldName)
			require.True(t, found, "Field '%s' not found in wrp.Message", fieldName)

			jsonTag := field.Tag.Get("json")
			list := strings.SplitN(jsonTag, ",", 2)
			require.NotEmpty(t, list[0], "Field '%s' does not have a JSON tag", fieldName)

			assert.Equal(t, list[0], expectedTag, "Constant for field '%s' does not match the JSON tag", fieldName)
		})
	}

	// Ensure all fields in wrp.Message are represented in the fieldMap
	for i := 0; i < msgType.NumField(); i++ {
		field := msgType.Field(i)
		if _, found := ignored[field.Name]; found {
			continue
		}

		assert.Contains(t, fieldMap, field.Name, "Field '%s' is not represented in the fieldMap", field.Name)
	}
}

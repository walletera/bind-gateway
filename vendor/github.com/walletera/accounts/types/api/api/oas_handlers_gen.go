// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/otelogen"
)

type codeRecorder struct {
	http.ResponseWriter
	status int
}

func (c *codeRecorder) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

// handleGetAccountRequest handles getAccount operation.
//
// Allows to retrieve an list of accounts based on search parameters.
//
// GET /accounts
func (s *Server) handleGetAccountRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	statusWriter := &codeRecorder{ResponseWriter: w}
	w = statusWriter
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("getAccount"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/accounts"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), GetAccountOperation,
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)

		attrSet := labeler.AttributeSet()
		attrs := attrSet.ToSlice()
		code := statusWriter.status
		if code != 0 {
			codeAttr := semconv.HTTPResponseStatusCode(code)
			attrs = append(attrs, codeAttr)
			span.SetAttributes(codeAttr)
		}
		attrOpt := metric.WithAttributes(attrs...)

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)

			// https://opentelemetry.io/docs/specs/semconv/http/http-spans/#status
			// Span Status MUST be left unset if HTTP status code was in the 1xx, 2xx or 3xx ranges,
			// unless there was another error (e.g., network error receiving the response body; or 3xx codes with
			// max redirects exceeded), in which case status MUST be set to Error.
			code := statusWriter.status
			if code >= 100 && code < 500 {
				span.SetStatus(codes.Error, stage)
			}

			attrSet := labeler.AttributeSet()
			attrs := attrSet.ToSlice()
			if code != 0 {
				attrs = append(attrs, semconv.HTTPResponseStatusCode(code))
			}

			s.errors.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: GetAccountOperation,
			ID:   "getAccount",
		}
	)
	params, err := decodeGetAccountParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetAccountRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    GetAccountOperation,
			OperationSummary: "Allows to retrieve an list of accounts based on search parameters",
			OperationID:      "getAccount",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "accountType",
					In:   "query",
				}: params.AccountType,
				{
					Name: "cvuAccountDetails",
					In:   "query",
				}: params.CvuAccountDetails,
				{
					Name: "dinopayAccountDetails",
					In:   "query",
				}: params.DinopayAccountDetails,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetAccountParams
			Response = GetAccountRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetAccountParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetAccount(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetAccount(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetAccountResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handlePostAccountRequest handles postAccount operation.
//
// Creates a new account.
//
// POST /accounts
func (s *Server) handlePostAccountRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	statusWriter := &codeRecorder{ResponseWriter: w}
	w = statusWriter
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("postAccount"),
		semconv.HTTPRequestMethodKey.String("POST"),
		semconv.HTTPRouteKey.String("/accounts"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), PostAccountOperation,
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)

		attrSet := labeler.AttributeSet()
		attrs := attrSet.ToSlice()
		code := statusWriter.status
		if code != 0 {
			codeAttr := semconv.HTTPResponseStatusCode(code)
			attrs = append(attrs, codeAttr)
			span.SetAttributes(codeAttr)
		}
		attrOpt := metric.WithAttributes(attrs...)

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)

			// https://opentelemetry.io/docs/specs/semconv/http/http-spans/#status
			// Span Status MUST be left unset if HTTP status code was in the 1xx, 2xx or 3xx ranges,
			// unless there was another error (e.g., network error receiving the response body; or 3xx codes with
			// max redirects exceeded), in which case status MUST be set to Error.
			code := statusWriter.status
			if code >= 100 && code < 500 {
				span.SetStatus(codes.Error, stage)
			}

			attrSet := labeler.AttributeSet()
			attrs := attrSet.ToSlice()
			if code != 0 {
				attrs = append(attrs, semconv.HTTPResponseStatusCode(code))
			}

			s.errors.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: PostAccountOperation,
			ID:   "postAccount",
		}
	)
	params, err := decodePostAccountParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	request, close, err := s.decodePostAccountRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response PostAccountRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    PostAccountOperation,
			OperationSummary: "Creates a new account",
			OperationID:      "postAccount",
			Body:             request,
			Params: middleware.Parameters{
				{
					Name: "X-Walletera-Correlation-Id",
					In:   "header",
				}: params.XWalleteraCorrelationID,
			},
			Raw: r,
		}

		type (
			Request  = *Account
			Params   = PostAccountParams
			Response = PostAccountRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackPostAccountParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.PostAccount(ctx, request, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.PostAccount(ctx, request, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodePostAccountResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

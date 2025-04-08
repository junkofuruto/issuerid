package issuerid

import (
	"context"
	"crypto/sha1"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type contextKey string

var (
	issuerIdKey = contextKey("CMG_MIDDLEWARE_ISSUER_ID")

	headerTrueClientIp = "True-Client-IP"
	headerForwardedFor = "X-Forwarded-For"
	headerRealIp       = "X-Real-IP"

	zeroUUID = "00000000-0000-0000-0000-000000000000"
)

func IssuerId(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if rip := realIP(r); rip != "" {
			uuid := generateUUIDFromString(rip)
			ctx = context.WithValue(ctx, issuerIdKey, contextKey(uuid.String()))
		} else {
			ctx = context.WithValue(ctx, issuerIdKey, zeroUUID)
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func GetIssuerID(ctx context.Context) string {
	if ctx == nil {
		return zeroUUID
	}

	if reqID, ok := ctx.Value(issuerIdKey).(string); ok {
		return reqID
	}

	return zeroUUID
}

func generateUUIDFromString(s string) uuid.UUID {
	h := sha1.New()
	h.Write([]byte(s))
	hash := h.Sum(nil)

	var u uuid.UUID
	copy(u[:], hash[:16])
	u[6] = (u[6] & 0x0f) | 0x50
	u[8] = (u[8] & 0x3f) | 0x80

	return u
}

func realIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get(headerTrueClientIp); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(headerForwardedFor); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(headerRealIp); xff != "" {
		ip, _, _ = strings.Cut(xff, ",")
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}

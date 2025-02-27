package felo

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"strings"
)

// Apple 服务器通知结构体
type ResponseBodyV2DecodedPayload struct {
	NotificationType string       `json:"notificationType"`
	Version          string       `json:"version"`
	Data             AppStoreData `json:"data"`
	SignedDate       int64        `json:"signedDate"`
	NotificationUUID string       `json:"notificationUUID"` // 幂等标识
}

// AppStoreData 结构体（存储交易数据）
type AppStoreData struct {
	AppAppleID     int64  `json:"appAppleId"`
	BundleID       string `json:"bundleId"`
	Environment    string `json:"environment"`
	JWSRenewalInfo string `json:"jwsRenewalInfo"`
	JWSTransaction string `json:"signedTransactionInfo"`
}

// JWT Payload 结构体（示例，可根据 Apple 文档扩展）
type AppStoreTransaction struct {
	TransactionID         string `json:"transactionId"`
	OriginalTransactionID string `json:"originalTransactionId"`
	ProductID             string `json:"productId"`
	PurchaseDate          int64  `json:"purchaseDate"`
	ExpiresDate           int64  `json:"expiresDate,omitempty"`
}

// `JWSTransactionDecodedPayload` 结构体（交易信息）
type JWSTransactionDecodedPayload struct {
	TransactionID         string `json:"transactionId"`
	OriginalTransactionID string `json:"originalTransactionId"`
	ProductID             string `json:"productId"`
	PurchaseDate          int64  `json:"purchaseDate"`
	ExpiresDate           *int64 `json:"expiresDate,omitempty"` // 订阅才会有
	RevocationDate        *int64 `json:"revocationDate,omitempty"`
	RevocationReason      *int   `json:"revocationReason,omitempty"`
	Quantity              int    `json:"quantity"`
	Type                  string `json:"type"`
	WebOrderLineItemID    string `json:"webOrderLineItemId"`
	SubscriptionGroupID   string `json:"subscriptionGroupIdentifier,omitempty"`
	OfferType             *int   `json:"offerType,omitempty"`
	OfferID               string `json:"offerId,omitempty"`
	InAppOwnershipType    string `json:"inAppOwnershipType"`
	SignedDate            int64  `json:"signedDate"`
	Environment           string `json:"environment"`
}

// 解析 JWT 并提取 Payload
func parseSignedPayload() {
	signedPayload := "eyJhbGciOiJFUzI1NiIsIng1YyI6WyJNSUlFTURDQ0E3YWdBd0lCQWdJUWZUbGZkMGZOdkZXdnpDMVlJQU5zWGpBS0JnZ3Foa2pPUFFRREF6QjFNVVF3UWdZRFZRUURERHRCY0hCc1pTQlhiM0pzWkhkcFpHVWdSR1YyWld4dmNHVnlJRkpsYkdGMGFXOXVjeUJEWlhKMGFXWnBZMkYwYVc5dUlFRjFkR2h2Y21sMGVURUxNQWtHQTFVRUN3d0NSell4RXpBUkJnTlZCQW9NQ2tGd2NHeGxJRWx1WXk0eEN6QUpCZ05WQkFZVEFsVlRNQjRYRFRJek1Ea3hNakU1TlRFMU0xb1hEVEkxTVRBeE1URTVOVEUxTWxvd2daSXhRREErQmdOVkJBTU1OMUJ5YjJRZ1JVTkRJRTFoWXlCQmNIQWdVM1J2Y21VZ1lXNWtJR2xVZFc1bGN5QlRkRzl5WlNCU1pXTmxhWEIwSUZOcFoyNXBibWN4TERBcUJnTlZCQXNNSTBGd2NHeGxJRmR2Y214a2QybGtaU0JFWlhabGJHOXdaWElnVW1Wc1lYUnBiMjV6TVJNd0VRWURWUVFLREFwQmNIQnNaU0JKYm1NdU1Rc3dDUVlEVlFRR0V3SlZVekJaTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCRUZFWWUvSnFUcXlRdi9kdFhrYXVESENTY1YxMjlGWVJWLzB4aUIyNG5DUWt6UWYzYXNISk9OUjVyMFJBMGFMdko0MzJoeTFTWk1vdXZ5ZnBtMjZqWFNqZ2dJSU1JSUNCREFNQmdOVkhSTUJBZjhFQWpBQU1COEdBMVVkSXdRWU1CYUFGRDh2bENOUjAxREptaWc5N2JCODVjK2xrR0taTUhBR0NDc0dBUVVGQndFQkJHUXdZakF0QmdnckJnRUZCUWN3QW9ZaGFIUjBjRG92TDJObGNuUnpMbUZ3Y0d4bExtTnZiUzkzZDJSeVp6WXVaR1Z5TURFR0NDc0dBUVVGQnpBQmhpVm9kSFJ3T2k4dmIyTnpjQzVoY0hCc1pTNWpiMjB2YjJOemNEQXpMWGQzWkhKbk5qQXlNSUlCSGdZRFZSMGdCSUlCRlRDQ0FSRXdnZ0VOQmdvcWhraUc5Mk5rQlFZQk1JSCtNSUhEQmdnckJnRUZCUWNDQWpDQnRneUJzMUpsYkdsaGJtTmxJRzl1SUhSb2FYTWdZMlZ5ZEdsbWFXTmhkR1VnWW5rZ1lXNTVJSEJoY25SNUlHRnpjM1Z0WlhNZ1lXTmpaWEIwWVc1alpTQnZaaUIwYUdVZ2RHaGxiaUJoY0hCc2FXTmhZbXhsSUhOMFlXNWtZWEprSUhSbGNtMXpJR0Z1WkNCamIyNWthWFJwYjI1eklHOW1JSFZ6WlN3Z1kyVnlkR2xtYVdOaGRHVWdjRzlzYVdONUlHRnVaQ0JqWlhKMGFXWnBZMkYwYVc5dUlIQnlZV04wYVdObElITjBZWFJsYldWdWRITXVNRFlHQ0NzR0FRVUZCd0lCRmlwb2RIUndPaTh2ZDNkM0xtRndjR3hsTG1OdmJTOWpaWEowYVdacFkyRjBaV0YxZEdodmNtbDBlUzh3SFFZRFZSME9CQllFRkFNczhQanM2VmhXR1FsekUyWk9FK0dYNE9vL01BNEdBMVVkRHdFQi93UUVBd0lIZ0RBUUJnb3Foa2lHOTJOa0Jnc0JCQUlGQURBS0JnZ3Foa2pPUFFRREF3Tm9BREJsQWpFQTh5Uk5kc2twNTA2REZkUExnaExMSndBdjVKOGhCR0xhSThERXhkY1BYK2FCS2pqTzhlVW85S3BmcGNOWVVZNVlBakFQWG1NWEVaTCtRMDJhZHJtbXNoTnh6M05uS20rb3VRd1U3dkJUbjBMdmxNN3ZwczJZc2xWVGFtUllMNGFTczVrPSIsIk1JSURGakNDQXB5Z0F3SUJBZ0lVSXNHaFJ3cDBjMm52VTRZU3ljYWZQVGp6Yk5jd0NnWUlLb1pJemowRUF3TXdaekViTUJrR0ExVUVBd3dTUVhCd2JHVWdVbTl2ZENCRFFTQXRJRWN6TVNZd0pBWURWUVFMREIxQmNIQnNaU0JEWlhKMGFXWnBZMkYwYVc5dUlFRjFkR2h2Y21sMGVURVRNQkVHQTFVRUNnd0tRWEJ3YkdVZ1NXNWpMakVMTUFrR0ExVUVCaE1DVlZNd0hoY05NakV3TXpFM01qQXpOekV3V2hjTk16WXdNekU1TURBd01EQXdXakIxTVVRd1FnWURWUVFERER0QmNIQnNaU0JYYjNKc1pIZHBaR1VnUkdWMlpXeHZjR1Z5SUZKbGJHRjBhVzl1Y3lCRFpYSjBhV1pwWTJGMGFXOXVJRUYxZEdodmNtbDBlVEVMTUFrR0ExVUVDd3dDUnpZeEV6QVJCZ05WQkFvTUNrRndjR3hsSUVsdVl5NHhDekFKQmdOVkJBWVRBbFZUTUhZd0VBWUhLb1pJemowQ0FRWUZLNEVFQUNJRFlnQUVic1FLQzk0UHJsV21aWG5YZ3R4emRWSkw4VDBTR1luZ0RSR3BuZ24zTjZQVDhKTUViN0ZEaTRiQm1QaENuWjMvc3E2UEYvY0djS1hXc0w1dk90ZVJoeUo0NXgzQVNQN2NPQithYW85MGZjcHhTdi9FWkZibmlBYk5nWkdoSWhwSW80SDZNSUgzTUJJR0ExVWRFd0VCL3dRSU1BWUJBZjhDQVFBd0h3WURWUjBqQkJnd0ZvQVV1N0Rlb1ZnemlKcWtpcG5ldnIzcnI5ckxKS3N3UmdZSUt3WUJCUVVIQVFFRU9qQTRNRFlHQ0NzR0FRVUZCekFCaGlwb2RIUndPaTh2YjJOemNDNWhjSEJzWlM1amIyMHZiMk56Y0RBekxXRndjR3hsY205dmRHTmhaek13TndZRFZSMGZCREF3TGpBc29DcWdLSVltYUhSMGNEb3ZMMk55YkM1aGNIQnNaUzVqYjIwdllYQndiR1Z5YjI5MFkyRm5NeTVqY213d0hRWURWUjBPQkJZRUZEOHZsQ05SMDFESm1pZzk3YkI4NWMrbGtHS1pNQTRHQTFVZER3RUIvd1FFQXdJQkJqQVFCZ29xaGtpRzkyTmtCZ0lCQkFJRkFEQUtCZ2dxaGtqT1BRUURBd05vQURCbEFqQkFYaFNxNUl5S29nTUNQdHc0OTBCYUI2NzdDYUVHSlh1ZlFCL0VxWkdkNkNTamlDdE9udU1UYlhWWG14eGN4ZmtDTVFEVFNQeGFyWlh2TnJreFUzVGtVTUkzM3l6dkZWVlJUNHd4V0pDOTk0T3NkY1o0K1JHTnNZRHlSNWdtZHIwbkRHZz0iLCJNSUlDUXpDQ0FjbWdBd0lCQWdJSUxjWDhpTkxGUzVVd0NnWUlLb1pJemowRUF3TXdaekViTUJrR0ExVUVBd3dTUVhCd2JHVWdVbTl2ZENCRFFTQXRJRWN6TVNZd0pBWURWUVFMREIxQmNIQnNaU0JEWlhKMGFXWnBZMkYwYVc5dUlFRjFkR2h2Y21sMGVURVRNQkVHQTFVRUNnd0tRWEJ3YkdVZ1NXNWpMakVMTUFrR0ExVUVCaE1DVlZNd0hoY05NVFF3TkRNd01UZ3hPVEEyV2hjTk16a3dORE13TVRneE9UQTJXakJuTVJzd0dRWURWUVFEREJKQmNIQnNaU0JTYjI5MElFTkJJQzBnUnpNeEpqQWtCZ05WQkFzTUhVRndjR3hsSUVObGNuUnBabWxqWVhScGIyNGdRWFYwYUc5eWFYUjVNUk13RVFZRFZRUUtEQXBCY0hCc1pTQkpibU11TVFzd0NRWURWUVFHRXdKVlV6QjJNQkFHQnlxR1NNNDlBZ0VHQlN1QkJBQWlBMklBQkpqcEx6MUFjcVR0a3lKeWdSTWMzUkNWOGNXalRuSGNGQmJaRHVXbUJTcDNaSHRmVGpqVHV4eEV0WC8xSDdZeVlsM0o2WVJiVHpCUEVWb0EvVmhZREtYMUR5eE5CMGNUZGRxWGw1ZHZNVnp0SzUxN0lEdll1VlRaWHBta09sRUtNYU5DTUVBd0hRWURWUjBPQkJZRUZMdXczcUZZTTRpYXBJcVozcjY5NjYvYXl5U3JNQThHQTFVZEV3RUIvd1FGTUFNQkFmOHdEZ1lEVlIwUEFRSC9CQVFEQWdFR01Bb0dDQ3FHU000OUJBTURBMmdBTUdVQ01RQ0Q2Y0hFRmw0YVhUUVkyZTN2OUd3T0FFWkx1Tit5UmhIRkQvM21lb3locG12T3dnUFVuUFdUeG5TNGF0K3FJeFVDTUcxbWloREsxQTNVVDgyTlF6NjBpbU9sTTI3amJkb1h0MlFmeUZNbStZaGlkRGtMRjF2TFVhZ002QmdENTZLeUtBPT0iXX0.eyJub3RpZmljYXRpb25UeXBlIjoiT05FX1RJTUVfQ0hBUkdFIiwibm90aWZpY2F0aW9uVVVJRCI6ImM0YWI5OTZlLTU3OTEtNDRmZi1hMjgzLWVkN2M3MmIyYmVkMyIsImRhdGEiOnsiYXBwQXBwbGVJZCI6MTY1OTM0NDA4OSwiYnVuZGxlSWQiOiJjb20uY3luZXQueGluZ3FpIiwiYnVuZGxlVmVyc2lvbiI6IjQxMCIsImVudmlyb25tZW50IjoiU2FuZGJveCIsInNpZ25lZFRyYW5zYWN0aW9uSW5mbyI6ImV5SmhiR2NpT2lKRlV6STFOaUlzSW5nMVl5STZXeUpOU1VsRlRVUkRRMEUzWVdkQmQwbENRV2RKVVdaVWJHWmtNR1pPZGtaWGRucERNVmxKUVU1eldHcEJTMEpuWjNGb2EycFBVRkZSUkVGNlFqRk5WVkYzVVdkWlJGWlJVVVJFUkhSQ1kwaENjMXBUUWxoaU0wcHpXa2hrY0ZwSFZXZFNSMVl5V2xkNGRtTkhWbmxKUmtwc1lrZEdNR0ZYT1hWamVVSkVXbGhLTUdGWFduQlpNa1l3WVZjNWRVbEZSakZrUjJoMlkyMXNNR1ZVUlV4TlFXdEhRVEZWUlVOM2QwTlNlbGw0UlhwQlVrSm5UbFpDUVc5TlEydEdkMk5IZUd4SlJXeDFXWGswZUVONlFVcENaMDVXUWtGWlZFRnNWbFJOUWpSWVJGUkplazFFYTNoTmFrVTFUbFJGTVUweGIxaEVWRWt4VFZSQmVFMVVSVFZPVkVVeFRXeHZkMmRhU1hoUlJFRXJRbWRPVmtKQlRVMU9NVUo1WWpKUloxSlZUa1JKUlRGb1dYbENRbU5JUVdkVk0xSjJZMjFWWjFsWE5XdEpSMnhWWkZjMWJHTjVRbFJrUnpsNVdsTkNVMXBYVG14aFdFSXdTVVpPY0ZveU5YQmliV040VEVSQmNVSm5UbFpDUVhOTlNUQkdkMk5IZUd4SlJtUjJZMjE0YTJReWJHdGFVMEpGV2xoYWJHSkhPWGRhV0VsblZXMVdjMWxZVW5CaU1qVjZUVkpOZDBWUldVUldVVkZMUkVGd1FtTklRbk5hVTBKS1ltMU5kVTFSYzNkRFVWbEVWbEZSUjBWM1NsWlZla0phVFVKTlIwSjVjVWRUVFRRNVFXZEZSME5EY1VkVFRUUTVRWGRGU0VFd1NVRkNSVVpGV1dVdlNuRlVjWGxSZGk5a2RGaHJZWFZFU0VOVFkxWXhNamxHV1ZKV0x6QjRhVUl5Tkc1RFVXdDZVV1l6WVhOSVNrOU9ValZ5TUZKQk1HRk1ka28wTXpKb2VURlRXazF2ZFhaNVpuQnRNalpxV0ZOcVoyZEpTVTFKU1VOQ1JFRk5RbWRPVmtoU1RVSkJaamhGUVdwQlFVMUNPRWRCTVZWa1NYZFJXVTFDWVVGR1JEaDJiRU5PVWpBeFJFcHRhV2M1TjJKQ09EVmpLMnhyUjB0YVRVaEJSME5EYzBkQlVWVkdRbmRGUWtKSFVYZFpha0YwUW1kbmNrSm5SVVpDVVdOM1FXOVphR0ZJVWpCalJHOTJUREpPYkdOdVVucE1iVVozWTBkNGJFeHRUblppVXprelpESlNlVnA2V1hWYVIxWjVUVVJGUjBORGMwZEJVVlZHUW5wQlFtaHBWbTlrU0ZKM1QyazRkbUl5VG5walF6Vm9ZMGhDYzFwVE5XcGlNakIyWWpKT2VtTkVRWHBNV0dReldraEtiazVxUVhsTlNVbENTR2RaUkZaU01HZENTVWxDUmxSRFEwRlNSWGRuWjBWT1FtZHZjV2hyYVVjNU1rNXJRbEZaUWsxSlNDdE5TVWhFUW1kbmNrSm5SVVpDVVdORFFXcERRblJuZVVKek1VcHNZa2RzYUdKdFRteEpSemwxU1VoU2IyRllUV2RaTWxaNVpFZHNiV0ZYVG1oa1IxVm5XVzVyWjFsWE5UVkpTRUpvWTI1U05VbEhSbnBqTTFaMFdsaE5aMWxYVG1wYVdFSXdXVmMxYWxwVFFuWmFhVUl3WVVkVloyUkhhR3hpYVVKb1kwaENjMkZYVG1oWmJYaHNTVWhPTUZsWE5XdFpXRXByU1VoU2JHTnRNWHBKUjBaMVdrTkNhbUl5Tld0aFdGSndZakkxZWtsSE9XMUpTRlo2V2xOM1oxa3lWbmxrUjJ4dFlWZE9hR1JIVldkalJ6bHpZVmRPTlVsSFJuVmFRMEpxV2xoS01HRlhXbkJaTWtZd1lWYzVkVWxJUW5sWlYwNHdZVmRPYkVsSVRqQlpXRkpzWWxkV2RXUklUWFZOUkZsSFEwTnpSMEZSVlVaQ2QwbENSbWx3YjJSSVVuZFBhVGgyWkROa00weHRSbmRqUjNoc1RHMU9kbUpUT1dwYVdFb3dZVmRhY0ZreVJqQmFWMFl4WkVkb2RtTnRiREJsVXpoM1NGRlpSRlpTTUU5Q1FsbEZSa0ZOY3poUWFuTTJWbWhYUjFGc2VrVXlXazlGSzBkWU5FOXZMMDFCTkVkQk1WVmtSSGRGUWk5M1VVVkJkMGxJWjBSQlVVSm5iM0ZvYTJsSE9USk9hMEpuYzBKQ1FVbEdRVVJCUzBKblozRm9hMnBQVUZGUlJFRjNUbTlCUkVKc1FXcEZRVGg1VWs1a2MydHdOVEEyUkVaa1VFeG5hRXhNU25kQmRqVktPR2hDUjB4aFNUaEVSWGhrWTFCWUsyRkNTMnBxVHpobFZXODVTM0JtY0dOT1dWVlpOVmxCYWtGUVdHMU5XRVZhVEN0Uk1ESmhaSEp0YlhOb1RuaDZNMDV1UzIwcmIzVlJkMVUzZGtKVWJqQk1kbXhOTjNad2N6SlpjMnhXVkdGdFVsbE1OR0ZUY3pWclBTSXNJazFKU1VSR2FrTkRRWEI1WjBGM1NVSkJaMGxWU1hOSGFGSjNjREJqTW01MlZUUlpVM2xqWVdaUVZHcDZZazVqZDBObldVbExiMXBKZW1vd1JVRjNUWGRhZWtWaVRVSnJSMEV4VlVWQmQzZFRVVmhDZDJKSFZXZFZiVGwyWkVOQ1JGRlRRWFJKUldONlRWTlpkMHBCV1VSV1VWRk1SRUl4UW1OSVFuTmFVMEpFV2xoS01HRlhXbkJaTWtZd1lWYzVkVWxGUmpGa1IyaDJZMjFzTUdWVVJWUk5Ra1ZIUVRGVlJVTm5kMHRSV0VKM1lrZFZaMU5YTldwTWFrVk1UVUZyUjBFeFZVVkNhRTFEVmxaTmQwaG9ZMDVOYWtWM1RYcEZNMDFxUVhwT2VrVjNWMmhqVGsxNldYZE5la1UxVFVSQmQwMUVRWGRYYWtJeFRWVlJkMUZuV1VSV1VWRkVSRVIwUW1OSVFuTmFVMEpZWWpOS2MxcElaSEJhUjFWblVrZFdNbHBYZUhaalIxWjVTVVpLYkdKSFJqQmhWemwxWTNsQ1JGcFlTakJoVjFwd1dUSkdNR0ZYT1hWSlJVWXhaRWRvZG1OdGJEQmxWRVZNVFVGclIwRXhWVVZEZDNkRFVucFplRVY2UVZKQ1owNVdRa0Z2VFVOclJuZGpSM2hzU1VWc2RWbDVOSGhEZWtGS1FtZE9Wa0pCV1ZSQmJGWlVUVWhaZDBWQldVaExiMXBKZW1vd1EwRlJXVVpMTkVWRlFVTkpSRmxuUVVWaWMxRkxRemswVUhKc1YyMWFXRzVZWjNSNGVtUldTa3c0VkRCVFIxbHVaMFJTUjNCdVoyNHpUalpRVkRoS1RVVmlOMFpFYVRSaVFtMVFhRU51V2pNdmMzRTJVRVl2WTBkalMxaFhjMHcxZGs5MFpWSm9lVW8wTlhnelFWTlFOMk5QUWl0aFlXODVNR1pqY0hoVGRpOUZXa1ppYm1sQllrNW5Xa2RvU1dod1NXODBTRFpOU1VnelRVSkpSMEV4VldSRmQwVkNMM2RSU1UxQldVSkJaamhEUVZGQmQwaDNXVVJXVWpCcVFrSm5kMFp2UVZWMU4wUmxiMVpuZW1sS2NXdHBjRzVsZG5JemNuSTVja3hLUzNOM1VtZFpTVXQzV1VKQ1VWVklRVkZGUlU5cVFUUk5SRmxIUTBOelIwRlJWVVpDZWtGQ2FHbHdiMlJJVW5kUGFUaDJZakpPZW1ORE5XaGpTRUp6V2xNMWFtSXlNSFppTWs1NlkwUkJla3hYUm5kalIzaHNZMjA1ZG1SSFRtaGFlazEzVG5kWlJGWlNNR1pDUkVGM1RHcEJjMjlEY1dkTFNWbHRZVWhTTUdORWIzWk1NazU1WWtNMWFHTklRbk5hVXpWcVlqSXdkbGxZUW5kaVIxWjVZakk1TUZreVJtNU5lVFZxWTIxM2QwaFJXVVJXVWpCUFFrSlpSVVpFT0hac1EwNVNNREZFU20xcFp6azNZa0k0TldNcmJHdEhTMXBOUVRSSFFURlZaRVIzUlVJdmQxRkZRWGRKUWtKcVFWRkNaMjl4YUd0cFJ6a3lUbXRDWjBsQ1FrRkpSa0ZFUVV0Q1oyZHhhR3RxVDFCUlVVUkJkMDV2UVVSQ2JFRnFRa0ZZYUZOeE5VbDVTMjluVFVOUWRIYzBPVEJDWVVJMk56ZERZVVZIU2xoMVpsRkNMMFZ4V2tka05rTlRhbWxEZEU5dWRVMVVZbGhXV0cxNGVHTjRabXREVFZGRVZGTlFlR0Z5V2xoMlRuSnJlRlV6Vkd0VlRVa3pNM2w2ZGtaV1ZsSlVOSGQ0VjBwRE9UazBUM05rWTFvMEsxSkhUbk5aUkhsU05XZHRaSEl3YmtSSFp6MGlMQ0pOU1VsRFVYcERRMEZqYldkQmQwbENRV2RKU1V4aldEaHBUa3hHVXpWVmQwTm5XVWxMYjFwSmVtb3dSVUYzVFhkYWVrVmlUVUpyUjBFeFZVVkJkM2RUVVZoQ2QySkhWV2RWYlRsMlpFTkNSRkZUUVhSSlJXTjZUVk5aZDBwQldVUldVVkZNUkVJeFFtTklRbk5hVTBKRVdsaEtNR0ZYV25CWk1rWXdZVmM1ZFVsRlJqRmtSMmgyWTIxc01HVlVSVlJOUWtWSFFURlZSVU5uZDB0UldFSjNZa2RWWjFOWE5XcE1ha1ZNVFVGclIwRXhWVVZDYUUxRFZsWk5kMGhvWTA1TlZGRjNUa1JOZDAxVVozaFBWRUV5VjJoalRrMTZhM2RPUkUxM1RWUm5lRTlVUVRKWGFrSnVUVkp6ZDBkUldVUldVVkZFUkVKS1FtTklRbk5hVTBKVFlqSTVNRWxGVGtKSlF6Qm5VbnBOZUVwcVFXdENaMDVXUWtGelRVaFZSbmRqUjNoc1NVVk9iR051VW5CYWJXeHFXVmhTY0dJeU5HZFJXRll3WVVjNWVXRllValZOVWsxM1JWRlpSRlpSVVV0RVFYQkNZMGhDYzFwVFFrcGliVTExVFZGemQwTlJXVVJXVVZGSFJYZEtWbFY2UWpKTlFrRkhRbmx4UjFOTk5EbEJaMFZIUWxOMVFrSkJRV2xCTWtsQlFrcHFjRXg2TVVGamNWUjBhM2xLZVdkU1RXTXpVa05XT0dOWGFsUnVTR05HUW1KYVJIVlhiVUpUY0ROYVNIUm1WR3BxVkhWNGVFVjBXQzh4U0RkWmVWbHNNMG8yV1ZKaVZIcENVRVZXYjBFdlZtaFpSRXRZTVVSNWVFNUNNR05VWkdSeFdHdzFaSFpOVm5wMFN6VXhOMGxFZGxsMVZsUmFXSEJ0YTA5c1JVdE5ZVTVEVFVWQmQwaFJXVVJXVWpCUFFrSlpSVVpNZFhjemNVWlpUVFJwWVhCSmNWb3pjalk1TmpZdllYbDVVM0pOUVRoSFFURlZaRVYzUlVJdmQxRkdUVUZOUWtGbU9IZEVaMWxFVmxJd1VFRlJTQzlDUVZGRVFXZEZSMDFCYjBkRFEzRkhVMDAwT1VKQlRVUkJNbWRCVFVkVlEwMVJRMFEyWTBoRlJtdzBZVmhVVVZreVpUTjJPVWQzVDBGRldreDFUaXQ1VW1oSVJrUXZNMjFsYjNsb2NHMTJUM2RuVUZWdVVGZFVlRzVUTkdGMEszRkplRlZEVFVjeGJXbG9SRXN4UVROVlZEZ3lUbEY2TmpCcGJVOXNUVEkzYW1Ka2IxaDBNbEZtZVVaTmJTdFphR2xrUkd0TVJqRjJURlZoWjAwMlFtZEVOVFpMZVV0QlBUMGlYWDAuZXlKMGNtRnVjMkZqZEdsdmJrbGtJam9pTWpBd01EQXdNRGcwT0RRME5USTBNaUlzSW05eWFXZHBibUZzVkhKaGJuTmhZM1JwYjI1SlpDSTZJakl3TURBd01EQTRORGcwTkRVeU5ESWlMQ0ppZFc1a2JHVkpaQ0k2SW1OdmJTNWplVzVsZEM1NGFXNW5jV2tpTENKd2NtOWtkV04wU1dRaU9pSmplVzVsZEM1amIyTnBMakYxYzJRaUxDSndkWEpqYUdGelpVUmhkR1VpT2pFM016ZzNOVFUwT1RRd01EQXNJbTl5YVdkcGJtRnNVSFZ5WTJoaGMyVkVZWFJsSWpveE56TTROelUxTkRrME1EQXdMQ0p4ZFdGdWRHbDBlU0k2TVN3aWRIbHdaU0k2SWtOdmJuTjFiV0ZpYkdVaUxDSnBia0Z3Y0U5M2JtVnljMmhwY0ZSNWNHVWlPaUpRVlZKRFNFRlRSVVFpTENKemFXZHVaV1JFWVhSbElqb3hOek00TnpVMU5UQTBORFF6TENKbGJuWnBjbTl1YldWdWRDSTZJbE5oYm1SaWIzZ2lMQ0owY21GdWMyRmpkR2x2YmxKbFlYTnZiaUk2SWxCVlVrTklRVk5GSWl3aWMzUnZjbVZtY205dWRDSTZJbFZUUVNJc0luTjBiM0psWm5KdmJuUkpaQ0k2SWpFME16UTBNU0lzSW5CeWFXTmxJam81T1RBc0ltTjFjbkpsYm1ONUlqb2lWVk5FSW4wLmNMTmNkN1ZnLWs3SlRNN3pPT0lwZVZLZW9aQVdFdGJ5U1pEaU9hempuZ2VJR1ZycEtxOEJNSkt4QTdIVnVsbkRWVG9zQ21jTGFueUZMTlB5MGloRUVBIn0sInZlcnNpb24iOiIyLjAiLCJzaWduZWREYXRlIjoxNzM4NzU1NTA0NDcxfQ.cYEtnVKLYKJIghiYdZF3A71Og_9TzYdIc0jtfA1TPqGunSlzrEbgVY0XbgNx4WRGs15shMUqxaBhIEjjeNSh7w"
	////res := &AppStoreTransaction{}
	//
	parts := strings.Split(signedPayload, ".")
	if len(parts) != 3 {
		logx.Errorf("invalid JWT format, signedPayload:%v", signedPayload)
		return
	}
	//
	//// 解析 JWT 但不验证签名
	//token, _, err := new(jwt.Parser).ParseUnverified(signedPayload, jwt.MapClaims{})
	//if err != nil {
	//	logx.Errorf("ParseUnverified failed, err:%v", err)
	//	return
	//}
	//
	//claims, ok := token.Claims.(jwt.MapClaims)
	//if !ok {
	//	logx.Errorf("Claims failed")
	//	return
	//}
	//
	//logx.Infof("claims :%v", claims)
	//
	////no := &AppStoreNotification{
	////	NotificationType: claims["notificationType"].(string),
	////	Version:          claims["version"].(string),
	////	Data:             claims["data"],
	////	SignedDate:       claims["signedDate"].(int64),
	////	NotificationUUID: claims["notificationUUID"].(string),
	////}
	//
	//payloadBytes, err := json.Marshal(claims)
	//if err != nil {
	//	logx.Errorf("Claims Marshal failed, err:%v", err)
	//	return
	//}
	//
	//var payload ResponseBodyV2DecodedPayload
	//err = json.Unmarshal(payloadBytes, &payload)
	//if err != nil {
	//	logx.Errorf("Claims Unmarshal failed, err:%v", err)
	//	return
	//}
	//
	////logx.Infof("no :%v", payload)
	//
	//// 解析 `JWSTransaction`
	//transaction, err := parseJWSTransaction(payload.Data.JWSTransaction)
	//if err != nil {
	//	logx.Errorf("Error parsing JWSTransaction: %v", err)
	//	return
	//}
	//
	//logx.Infof("transaction :%v", transaction)

	claims, err := parseJWT(signedPayload)
	if err != nil {
		logx.Errorf("parse jwt failed, err:%v", err)
		return
	}

	dataJSON, err := json.Marshal(claims["data"])
	if err != nil {
		logx.Errorf("json marshal failed, err:%v", err)
		return
	}
	var appStoreData AppStoreData
	if err = json.Unmarshal(dataJSON, &appStoreData); err != nil {
		logx.Errorf("json unmarshal failed, err:%v", err)
		return
	}

	transactionClaims, err := parseJWT(appStoreData.JWSTransaction)
	if err != nil {
		logx.Errorf("parseJWT JWSTransaction failed, err:%v", err)
		return
	}

	// 解析 `JWSTransactionDecodedPayload`
	transactionJSON, _ := json.Marshal(transactionClaims)
	var transaction JWSTransactionDecodedPayload
	if err = json.Unmarshal(transactionJSON, &transaction); err != nil {
		logx.Errorf("json unmarshal failed, err:%v", err)
		return
	}

	logx.Infof("Received Transaction: %+v", transaction)
}

// 解析 JWSTransaction（交易信息）
func parseJWSTransaction(jwsTransaction string) (*JWSTransactionDecodedPayload, error) {
	parts := strings.Split(jwsTransaction, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	// 解析 JWT 但不验证签名
	token, _, err := new(jwt.Parser).ParseUnverified(jwsTransaction, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid JWT claims format")
	}

	// 解析交易数据
	transaction := &JWSTransactionDecodedPayload{}
	transactionBytes, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction claims: %w", err)
	}

	err = json.Unmarshal(transactionBytes, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return transaction, nil
}

// 解析 JWT（验证签名并解析 `payload`）
func parseJWT(signedJWT string) (jwt.MapClaims, error) {
	parts := strings.Split(signedJWT, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	// 解析 Header 以获取 `kid`
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %v", err)
	}

	var header map[string]interface{}
	if err = json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("failed to parse header: %v", err)
	}

	// 获取 Apple 公钥
	publicKey, err := getApplePublicKey(header["kid"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to get Apple public key: %w", err)
	}

	// 解析并验证 JWT
	token, err := jwt.Parse(signedJWT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid JWT claims")
	}

	return claims, nil
}

// 从 Apple 服务器获取 JWK 公钥
func getApplePublicKey(kid string) (*ecdsa.PublicKey, error) {
	resp, err := http.Get(appleJWKURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get Apple JWK: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWK response: %w", err)
	}

	var jwk AppleJWK
	err = json.Unmarshal(body, &jwk)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWK: %w", err)
	}

	for _, key := range jwk.Keys {
		if key.Kid == kid {
			// 获取 x5c（PEM 编码的公钥）
			certDER, err := base64.StdEncoding.DecodeString(key.X5c[0])
			if err != nil {
				return nil, fmt.Errorf("failed to decode x5c: %w", err)
			}

			cert, err := x509.ParseCertificate(certDER)
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate: %w", err)
			}

			// 解析 ECDSA 公钥
			pubKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
			if !ok {
				return nil, fmt.Errorf("failed to parse ECDSA public key")
			}

			return pubKey, nil
		}
	}

	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

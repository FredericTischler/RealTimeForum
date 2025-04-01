package middleware

import (
	"net/http"
	"sync"
	"time"

	"forum/services"
)

// userLimiter stocke le nombre de requêtes et le moment du dernier reset
type userLimiter struct {
	lastReset time.Time
	count     int
}

var (
	// Définir le nombre maximum de requêtes autorisées par intervalle
	// limitPerInterval = 100
	// Intervalle de temps pour le reset des compteurs (ici 1 minute)
	interval = time.Minute

	// Map pour stocker les compteurs par clé (UUID de l'utilisateur ou IP)
	mu           sync.Mutex
	userLimiters = make(map[string]*userLimiter)
)

// SessionRateLimiter est un middleware qui limite le nombre de requêtes
// par utilisateur (basé sur la session) ou par IP si non authentifié.
func SessionRateLimiter(sessionService *services.SessionService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key string

		// Essayer d'extraire le cookie "session_token"
		cookie, err := r.Cookie("session_token")
		if err == nil {
			// Si le cookie est présent, vérifier la session associée
			session, err := sessionService.GetSessionByToken(cookie.Value)
			if err == nil && session != nil {
				// Utiliser l'UUID de l'utilisateur comme clé
				key = session.UserId.String()
			}
		}

		// Si aucune session n'est trouvée, on utilise l'adresse IP comme fallback
		if key == "" {
			key = r.RemoteAddr
		}

		mu.Lock()
		limiter, exists := userLimiters[key]
		now := time.Now()
		// Si le compteur n'existe pas ou que l'intervalle est dépassé, on le réinitialise
		if !exists || now.Sub(limiter.lastReset) > interval {
			limiter = &userLimiter{
				lastReset: now,
				count:     0,
			}
			userLimiters[key] = limiter
		}

		// Si le nombre de requêtes dépasse la limite, on retourne une erreur HTTP 429
		// if limiter.count >= limitPerInterval {
		// 	mu.Unlock()
		// 	http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		// 	return
		// }

		// Sinon, on incrémente le compteur et on passe au handler suivant
		limiter.count++
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

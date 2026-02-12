/**
 * Copyright 2025 Saber authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

package sbnet

import (
	"os-artificer/saber/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Option configures a Server.
type Option func(*serverConfig)

type serverConfig struct {
	logger          logger.Logger
	authMiddlewares []gin.HandlerFunc
	registrars      []APIRegistrar
	routes          []Route
	requestLogging  bool
}

// WithLogger injects a logger for request and lifecycle logging. If not set,
// the package-level logger (pkg/logger) is used.
func WithLogger(l logger.Logger) Option {
	return func(c *serverConfig) {
		c.logger = l
	}
}

// WithAuthMiddleware adds one or more middleware handlers for authenticated
// routes. These are applied to groups returned by Server.AuthGroup.
func WithAuthMiddleware(mw ...gin.HandlerFunc) Option {
	return func(c *serverConfig) {
		c.authMiddlewares = append(c.authMiddlewares, mw...)
	}
}

// WithRegistrars registers RESTful API modules. Each registrar's Register(s) is
// called with the Server so it can use s.Engine() and s.AuthGroup(path).
func WithRegistrars(registrars ...APIRegistrar) Option {
	return func(c *serverConfig) {
		c.registrars = append(c.registrars, registrars...)
	}
}

// WithRoutes registers routes at construction. Each route binds method, path,
// and handler; use Get/Post/etc. for public routes and AuthGet/AuthPost/etc. for
// routes under an auth group. No need to implement APIRegistrar.
func WithRoutes(routes ...Route) Option {
	return func(c *serverConfig) {
		c.routes = append(c.routes, routes...)
	}
}

// WithRequestLogging enables or disables HTTP request logging middleware.
// Default is true.
func WithRequestLogging(enable bool) Option {
	return func(c *serverConfig) {
		c.requestLogging = enable
	}
}

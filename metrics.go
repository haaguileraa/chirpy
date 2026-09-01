package main

import (
	"fmt"
)

const metricsTemplate = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func getFormattedMetrics(visits int) string {
	return fmt.Sprintf(metricsTemplate, visits)
} 

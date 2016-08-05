# grbinder
Gin route binder

# usage
  ```go
  type statusHandler struct {
  }
  
  func (h *statusHandler) GET(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }
  
  router := gin.Default()
  grbinder.BindVerb(router.Group("/status"), &statusHandler{})
  
  ```

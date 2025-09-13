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

With lock

  ```go
  type entityHandler struct {
  }

  func (h *entityHandler) CreateHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }
  
  func (h *entityHandler) TakeHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }

  func (h *entityHandler) UpdateHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }

  func (h *entityHandler) DeleteHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }
  
  router := gin.Default()
  grbinder.CRUD(
    router.Group("/entity"),
    &entityHandler{},
    grbinder.WithEntityLockEnable(true),
    grbinder.WithEntityLockName("entity"),
  )
  ```

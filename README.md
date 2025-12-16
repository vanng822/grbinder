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
  func init() {
    grbinder.InitDefaultLocker()
  }

  type calendarHandler struct {
  }

  func (h *calendarHandler) CreateHandler(c *gin.Context) {
    // create calendar
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }
  
  func (h *calendarHandler) TakeHandler(c *gin.Context) {
    // serve calendar
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }

  func (h *calendarHandler) UpdateHandler(c *gin.Context) {
    // update calendar required lock
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }

  func (h *calendarHandler) DeleteHandler(c *gin.Context) {
    // delete calendar required lock
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }
  
  router := gin.Default()
  grbinder.CRUD(
    router.Group("/calendar"),
    &calendarHandler{},
    grbinder.WithEntityLockEnable(true),
    grbinder.WithEntityLockName("calendar"),
  )

  type calendarSubscribeHandler struct {
  }
  
  func (h *calendarSubscribeHandler) POST(c *gin.Context) {
    // Lock when subscribing calendar
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }

  func (h *calendarSubscribeHandler) DELETE(c *gin.Context) {
    // Lock when unsubscribing calendar
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  }

  grbinder.BindVerb(
    router.Group("/calendar/:id/subscribe"),
    &calendarSubscribeHandler{},
    grbinder.WithEntityLockEnable(true),
    grbinder.WithEntityLockName("calendar"),
  )
  ```

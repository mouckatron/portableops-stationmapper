package main

import (
  "fmt"
  
  "github.com/gin-gonic/gin"
)

func tileRouterPaths(r *gin.Engine){

  r.GET("/tile/:z/:x/:y", func(c *gin.Context){

    z := c.Param("z")
    x := c.Param("x")
    y := c.Param("y")
    
    c.Header("Content-Type", "image/png")
    filepath := fmt.Sprintf("%s/%s/%s/%s.png", rootpath, z, x, y)
    print("Looking for ")
    print(filepath)
    c.File(filepath)
  })
  
}

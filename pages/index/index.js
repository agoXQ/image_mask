
// index.js
Page({
    data:{
        imageSrc:"",
        imageDeal:"",
        imgpasswd:"",
        operation:"decrypt",
        id:""
    },
    operateChange:function(e){
        console.log(e)
        if (e.detail.value){
            this.setData({
                operation:"encrypt"   
            })
            console.log("set operation encrypt")
        }else if(e.detail.value===false){
            this.setData({
                operation:"decrypt"
            })    
            console.log("set operation decrypt")
        }else{
            console.log("error occur in set operation")
        }
    },
    keyInput: function(e){
        this.setData({
            imgpasswd:e.detail.value
            })
        console.log(this.data.imgpasswd)
    },
    chooseImage: function(event){
        console.log(event)
        this.setData({
            imageSrc:"",
            imageDeal:""
        })
        var that = this
        wx.chooseMedia({
            count:1,
            mediaType:["image","video"],
            sizeType:["original","compressed"],
            sourceType:["album","camera"],
            success(res){
                const tempFilePaths = res.tempFiles[0].tempFilePath
                that.setData({
                    imageSrc:tempFilePaths
                })
                console.log(tempFilePaths)
            },
        })
    },
    sendImage: function(event){
        console.log(event)
        var that = this
        const tempFilePaths = this.data.imageSrc
        // console.log(tempFilePaths==="")
        var id = Date.now()
        that.setData({
            id:id
        })
        var passwd = this.data.imgpasswd
        var oprate = this.data.operation
        if (tempFilePaths==""){
            console.log("wqdoja")
            wx.showToast({
              title: '图片未选择',
              icon:"error",
              duration:2000,
              success:function(){
                  console.log("toast success")
              },
              fail(err){
                  console.error(err)
              }
            })
            return
        }
        // console.log(tempFilePaths)
        
        else if (passwd==="" || isNaN(Number(passwd))){
            wx.showToast({
              title: '请设置混淆密码',
              icon:"error",
              duration:2000,
              success(){
                  console.log(passwd==="")
                  console.log(isNaN(Number(passwd)))
              }
            })
            return
        }
        that.setData({
            imageDeal:""
        })
        var url = "http://192.168.116.251:8888/image/mask/"+id+"/"+oprate+"/"+passwd
        console.log("正在发送图片。。。",passwd)
        console.log(url)
        wx.uploadFile({
            filePath: tempFilePaths,
            url: url,
            name:"file",
            success:function(ress){
              var data = ress.data
              try{
                const response = JSON.parse(data)
                var filename = response.file
                var downurl = "http://192.168.116.251:8888/image/download/"+id+"/"+filename
                that.downloadFile(downurl,id+"_"+filename)
                console.log(response)
                }catch(error){
                    console.log(error)
                }
            },
            fail:function(res){
                  console.log(res)
                  console.log("图片失败")
            }
        })
        console.log("发送图片完成！")
    },
    // sendImage: function(event){
    //     console.log(event)
    //     var that = this
    //     wx.chooseMedia({
    //         count:1,
    //         mediaType:["image","video"],
    //         sizeType:["original","compressed"],
    //         sourceType:["album","camera"],
    //         success(res){
    //             const tempFilePaths = res.tempFiles[0].tempFilePath
    //             that.setData({
    //                 imageSrc:tempFilePaths
    //             })
    //             console.log(tempFilePaths)
    //             var url = "http://192.168.116.251:8888/image/mask/114514/encrypt"
    //             console.log("正在发送图片。。。")
    //             wx.uploadFile({
    //                 filePath: tempFilePaths,
    //                 url: url,
    //                 name:"file",
    //                 success:function(ress){
    //                   var data = ress.data
    //                   try{
    //                     const response = JSON.parse(data)
    //                     var filename = response.file
    //                     var downurl = "http://192.168.116.251:8888/image/download/114514/"+filename
    //                     that.downloadFile(downurl,filename)
    //                     console.log(response)
    //                     }catch(error){
    //                         console.log(error)
    //                     }
    //                 },
    //                 fail:function(res){
    //                       console.log(res)
    //                       console.log("图片失败")
    //                 }
    //             })
    //             console.log("发送图片完成！")
    //         },
    //         fail(){
    //             console.log("发送图片失败！")
    //             return
    //         }
    //     })
      
    // },
downloadFile: function(url,filename){
    // tempFilePath=
    const tempfile = wx.env.USER_DATA_PATH+'/'+filename
    var that =this
    wx.downloadFile({
      url: url,
      filePath:tempfile,
      success(res){
          console.log(res)
          if (res.statusCode===200){
            that.setData({
                imageDeal:tempfile
            })
            console.log(that.data.imageSrc,that.data.imageDeal)
            console.log("image download success")
          }
      },
      fail(error){
          console.log("image download error")
          console.log(error)
      }
    })
},
    saveImageToLocal:function(){
        const imageUrl = this.data.imageDeal
        wx.downloadFile({
          url: imageUrl,
          success(res){
            if (res.statusCode===200){
                const tempFailPath = res.tempFilePath
                wx.saveImageToPhotosAlbum({
                  filePath: tempFailPath,
                  success(){
                      wx.showToast({
                        title: '图片已保存',
                        icon:"success",
                        duration:2000
                      })
                  },
                  fail(err){
                      console.error("图片保存失败：",err)
                      wx.showToast({
                        title: '图片保存失败',
                        icon:"none",
                        duration:2000
                      })
                  }
                })
            }else{
                console.error("下载图片失败",res)
                wx.showToast({
                  title: '下载失败',
                  icon:"none",
                  duration:2000
                })
            }
          }
        })
    }
})
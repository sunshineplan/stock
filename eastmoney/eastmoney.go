package eastmoney

const ssePattern = `000[0-1]\d{2}|(51[0-358]|60[0-3]|688)\d{3}`
const szsePattern = `(00[0-3]|159|300|399)\d{3}`

const api = "http://push2his.eastmoney.com/api/qt/stock/trends2/get?iscr=0&fields1=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13&fields2=f51,f52,f53,f54,f55,f56,f57,f58&secid=1.600519"
const suggestAPI = "http://smartbox.gtimg.cn/s3/?t=%s&q=%s"

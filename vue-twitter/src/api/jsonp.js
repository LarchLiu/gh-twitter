import originJsonp from "jsonp"

export default function jsonp(url, data, option){   //url是不带参数的url路径，data是对象类型参数
  url += (url.indexOf('?') < 0 ? "?" : "&") + param(data);      //判断url中有没有?如果没有就拼接"?",如果有就拼接'$'然后在拼接参数，最终拼接成get传参的方式发送请求
  return new Promise((resolve,reject)=>{
    originJsonp(url,option,(err,data)=>{
      if(!err){
        resolve(data)
      }else{
        reject(err)
      }
    })
  })
}

export function param(data){   //将对象参数data里的键值对参数拼接在路径url上
  let url=""
  for(let k in data){
    let val = data[k] !== undefined ? data[k] : '';
    url +=`&${k}=${encodeURIComponent(val)}`   //url上的参数值可能为中文，所有需要encodeURIComponent编码
  }
  return url?url.substring(1):''    //如果url有数据，则去掉第一个$符号
}


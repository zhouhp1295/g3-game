(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-75d55c9b"],{"0f7d":function(e,t,r){},7803:function(e,t,r){"use strict";r.r(t);var s=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"register"},[r("el-form",{ref:"registerForm",staticClass:"register-form",attrs:{model:e.registerForm,rules:e.registerRules}},[r("h3",{staticClass:"title"},[e._v("g3-cms 管理后台")]),r("el-form-item",{attrs:{prop:"username"}},[r("el-input",{attrs:{type:"text","auto-complete":"off",placeholder:"账号"},model:{value:e.registerForm.username,callback:function(t){e.$set(e.registerForm,"username",t)},expression:"registerForm.username"}},[r("svg-icon",{staticClass:"el-input__icon input-icon",attrs:{slot:"prefix","icon-class":"user"},slot:"prefix"})],1)],1),r("el-form-item",{attrs:{prop:"password"}},[r("el-input",{attrs:{type:"password","auto-complete":"off",placeholder:"密码"},nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.handleRegister(t)}},model:{value:e.registerForm.password,callback:function(t){e.$set(e.registerForm,"password",t)},expression:"registerForm.password"}},[r("svg-icon",{staticClass:"el-input__icon input-icon",attrs:{slot:"prefix","icon-class":"password"},slot:"prefix"})],1)],1),r("el-form-item",{attrs:{prop:"confirmPassword"}},[r("el-input",{attrs:{type:"password","auto-complete":"off",placeholder:"确认密码"},nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.handleRegister(t)}},model:{value:e.registerForm.confirmPassword,callback:function(t){e.$set(e.registerForm,"confirmPassword",t)},expression:"registerForm.confirmPassword"}},[r("svg-icon",{staticClass:"el-input__icon input-icon",attrs:{slot:"prefix","icon-class":"password"},slot:"prefix"})],1)],1),e.captchaOnOff?r("el-form-item",{attrs:{prop:"code"}},[r("el-input",{staticStyle:{width:"63%"},attrs:{"auto-complete":"off",placeholder:"验证码"},nativeOn:{keyup:function(t){return!t.type.indexOf("key")&&e._k(t.keyCode,"enter",13,t.key,"Enter")?null:e.handleRegister(t)}},model:{value:e.registerForm.code,callback:function(t){e.$set(e.registerForm,"code",t)},expression:"registerForm.code"}},[r("svg-icon",{staticClass:"el-input__icon input-icon",attrs:{slot:"prefix","icon-class":"validCode"},slot:"prefix"})],1),r("div",{staticClass:"register-code"},[r("img",{staticClass:"register-code-img",attrs:{src:e.codeUrl},on:{click:e.getCode}})])],1):e._e(),r("el-form-item",{staticStyle:{width:"100%"}},[r("el-button",{staticStyle:{width:"100%"},attrs:{loading:e.loading,size:"medium",type:"primary"},nativeOn:{click:function(t){return t.preventDefault(),e.handleRegister(t)}}},[e.loading?r("span",[e._v("注 册 中...")]):r("span",[e._v("注 册")])]),r("div",{staticStyle:{float:"right"}},[r("router-link",{staticClass:"link-type",attrs:{to:"/login"}},[e._v("使用已有账户登录")])],1)],1)],1),e._m(0)],1)},i=[function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"el-register-footer"},[r("span",[e._v("Copyright © 2018-2021 ruoyi.vip All Rights Reserved.")])])}],o=(r("d9e2"),r("7ded")),n={name:"Register",data:function(){var e=this,t=function(t,r,s){e.registerForm.password!==r?s(new Error("两次输入的密码不一致")):s()};return{codeUrl:"",registerForm:{username:"",password:"",confirmPassword:"",code:"",uuid:""},registerRules:{username:[{required:!0,trigger:"blur",message:"请输入您的账号"},{min:2,max:20,message:"用户账号长度必须介于 2 和 20 之间",trigger:"blur"}],password:[{required:!0,trigger:"blur",message:"请输入您的密码"},{min:5,max:20,message:"用户密码长度必须介于 5 和 20 之间",trigger:"blur"}],confirmPassword:[{required:!0,trigger:"blur",message:"请再次输入您的密码"},{required:!0,validator:t,trigger:"blur"}],code:[{required:!0,trigger:"change",message:"请输入验证码"}]},loading:!1,captchaOnOff:!0}},created:function(){this.getCode()},methods:{getCode:function(){var e=this;Object(o["a"])().then((function(t){e.captchaOnOff=void 0===t.captchaOnOff||t.captchaOnOff,e.captchaOnOff&&(e.codeUrl="data:image/gif;base64,"+t.img,e.registerForm.uuid=t.uuid)}))},handleRegister:function(){var e=this;this.$refs.registerForm.validate((function(t){t&&(e.loading=!0,Object(o["e"])(e.registerForm).then((function(t){var r=e.registerForm.username;e.$alert("<font color='red'>恭喜你，您的账号 "+r+" 注册成功！</font>","系统提示",{dangerouslyUseHTMLString:!0,type:"success"}).then((function(){e.$router.push("/login")})).catch((function(){}))})).catch((function(){e.loading=!1,e.captchaOnOff&&e.getCode()})))}))}}},a=n,c=(r("e1bc"),r("2877")),l=Object(c["a"])(a,s,i,!1,null,null,null);t["default"]=l.exports},e1bc:function(e,t,r){"use strict";r("0f7d")}}]);
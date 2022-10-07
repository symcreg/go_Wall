function returnBack(){
    window.location.href="../index.html"
}
function gotoLogin(){
    window.location.href="./login/login.html"
}
function gotoRegister(){
    console.log("1");
    window.location.href="./register/register.html"
}
function gotoWall(){
    window.location.href="./wall/wall.html"
}
function gotoAdd(){
    var username=window.location.href.split("=")[1];
    window.location.href="../user-manager/add/add.html?username="+ username;
}
function gotoManage(){
    var username=window.location.href.split("=")[1];
    window.location.href="../../user-manager/manage/manage.html?username="+ username;
}
function reload(){
 location.reload();   
}
function quit(){
    var isQuit=window.confirm("是否退出登录");
    if (isQuit)
        window.location.href="../index.html";
}
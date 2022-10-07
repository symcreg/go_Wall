function login(){
    var data={};
    const url="http://localhost:8080/login";
    data.username="hxy";
    data.password="xx";
    console.log(data);
    $.ajax({
        url:url,
        data:JSON.stringify(data),
        type:"POST",
        contentType:'application/json',
        success:(res)=>{
            console.log(res);
            window.location.href="../wall/login.html?username="+data.username;
        }
    })
}
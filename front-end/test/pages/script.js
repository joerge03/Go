 // TEMPLATE LITTERAL ${}

 var test = "test"

 const unknownLetter =  `asdfasdf ${test}` 

 // result = asdfasdf test
var testButton = document.getElementById("brokerBtn")
var authBroker = document.getElementById("authBrokerBtn")
var loggerButton = document.getElementById("loggerBtn")
var mailerButton = document.getElementById("mailerBtn")
var output =  document.getElementById("output")
var payload = document.getElementById("payload")
var received = document.getElementById("received")

const test = se

mailerButton.addEventListener("click", async ()=>{
 console.log("mailer button clicked")

 try {          
   const payloadBody = {
    action: "mail",
     mail: {
       from: "test@gmail.com",
       fromName: "michael v the jackson 92",
       to: "joergemichaels@gmail.com",
       subject: "test",
       message: "panget"
     }
   }

   const pay = {
     method: "POST",
     headers: {"Content-Type":"application/json"},
     body: JSON.stringify(payloadBody),
   }

   const mailURL = "http://localhost:8085/handle"
   let res = await fetch(mailURL, pay)
   if (!res.ok) {
     throw Error(`There's something wrong getting response from: ${mailURL}`)
   }

   const resData = await res.json()
   console.log(res)
   output.innerHTML = JSON.stringify(resData)
   payload.innerHTML = JSON.stringify(payloadBody)         
   received.innerHTML = JSON.stringify()
 } catch (err) {
    output.innerHTML = `<br> <h1> Error: </h1> ${error}`  
 }
})

loggerButton.addEventListener("click",async ()=>  {
 
 console.log("logger button clicked")

 try {
   const pay = {
     action: "log",
     log: {
       name: "batumbakal",
       data: {"testing": "asdf" }
     }
   }
   const data = {
     method: "POST",
     headers: {
       "Content-Type": "application/json" 
     },
     body: JSON.stringify(pay) 
   }
   let res  = await fetch("http://localhost:8085/handle", data)
   if (!res.ok){
   throw Error("there is something wrong fetching handle")
 }

   const resData = await res.json()
   if (resData.error){
     console.log("there is something wrong fetching log")            
   }
   output.innerHTML = JSON.stringify(resData, undefined, 4)
   payload.innerHTML = JSON.stringify(pay, undefined, 4)
   received.innerHTML = JSON.stringify(resData)     
 } catch (err) {
     output.innerHTML = `<br> <h1> Error: </h1> ${error}`       
 }       
})

authBroker.addEventListener("click", async function () {
 console.log("clicked 32")
 const pay = { 
     action: "auth",
     auth: {
         email: "admin@example.com",
         password: "verysecret"
     }
   }
 const body = {
   method: "POST",
   headers: {"Content-Type": "application/json"},          
   body: JSON.stringify(pay)
 }
 console.log(JSON.stringify(pay))

 try {          
   const res = await fetch("http:\/\/localhost:8085/handle", body)    
   
   if (!res.ok){
     throw Error("there is something wrong fetching handle")
   }
   const data = await res.json()          

   payload.innerHTML = JSON.stringify(pay, undefined, 4)
   received.innerHTML = JSON.stringify(data)          
   if (data.error){
     console.log(data.error)
     output.innerHTML = `<br/> <h1> error: ${data.error}</h1>`           
   }
 } catch (error) {
     console.log(error)
     output.innerHTML = `<br> <h1> Error: </h1> ${error}`
 }
})


testButton.addEventListener("click", function(){
 const body= {
   method: 'POST',
   headers: {"Content-Type": 'application/json'}
 }

 testButton.innerText = "testing"

 fetch("http:\/\/localhost:8085", body)
 .then((response)=> {
   console.log(response)
   return response.json()    
 })
 .then((data)=>{
   payload.innerHTML = "empty post request"
   received.innerHTML = JSON.stringify(data, undefined, 4)
   if (data){
     // console.log(data.error)
     console.log(data)
   }else{
     console.log(data)
     output.innerHTML += `<br> <strong>Response from the server </strong>: ${data.message}` 
   }
 }).catch((err)=> {
   console.log(err)
     output.innerHTML = `<br> <h1> Error: </h1> ${err}`
 })
})
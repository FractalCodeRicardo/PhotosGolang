package main

import (
    "strconv"
    "bytes"
    "encoding/json"
    "fmt"
    "io"
	"io/ioutil"
    "net/http"
)

//Constantes
const API_URL = "https://jsonplaceholder.typicode.com/photos"

//hice un incoming webhook  de prueba no estoy seguro cuánto duren
const API_SLACk = "https://hooks.slack.com/services/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"

//Opciones del menú
const (
    OpcionListar = 1
    OpcionCrear = 2
    OpcionModificar  = 3
    OpcionConsultar =4
    OpcionEliminar = 5
    OpcionSalir = 6
)

// Estructura para la imagen
type Image struct{
    AlbumId int  `json:albumId` 
    Id int  `json:id` 
    Title string  `json:title` 
    Url string  `json:url` 
    ThumbnailUrl string  `json:thumbnailUrl` 
}

//Método main, contiene el ciclo del menú
func main() {

    opcion := 1

    for {
        opcion = Menu()

        //Esta variable sirve para almacenar el error de las peticiones
        var error error

        //Esta variable sirve para almacenar el mensaje de éxito de las peticiones
        var exito string


        /*Opciones del menú*/
        if opcion == OpcionListar {
            exito, error = ListarImagenes()
        }

        if opcion == OpcionCrear {
            exito, error = CrearImagen()
        }

        if opcion == OpcionModificar {
            exito, error = ModificarImagen()
        }

        if opcion == OpcionConsultar {
            exito, error = ConsultarImagen()
        }

        if opcion == OpcionEliminar {
            exito, error = BorrarImagen()
        }

        if opcion == OpcionSalir {
            break
        }

        //Si el usuario oprime una opción inválida
        if !EsOpcionValida(opcion){
            fmt.Println("Opción no válida")            
        }


        if error != nil {
            fmt.Println(error)
        }

        
        if exito != "" {
            fmt.Println(exito)
        }


        //Le pregunta al usuario si desea hacer otra operación
        if !OtraOperacion(){
            break;
        }
    }

}

//Devuelve verdadero si el usuario indica que quiere hacer otra operación
func OtraOperacion() bool {
    var input string
    fmt.Println("¿Quiere realizar otra operación? (S/N)")
    fmt.Scanln(&input)
    return input == "S" || input == "s"

}

//Devuelve verdadero si opción está entra 0 a 7
func EsOpcionValida(opcion int) bool{
    return opcion>0 && opcion<7
}

//Despliega el menú y pide al usuario que ingrese una opción
func Menu() int {
    fmt.Println("Elige una opción")

    fmt.Println("1. Listar imagenes")
    fmt.Println("2. Crear imagen")
    fmt.Println("3. Modificar imagen")
    fmt.Println("4. Consultar una imagen:")
    fmt.Println("5. Eliminar imagen")
    fmt.Println("6. Salir")

    fmt.Print("\nOpción:")
    var input string
    fmt.Scanln(&input)

    res, err := strconv.Atoi(input);

    //En caso de que el usuario ponga una opción inválida se retorna -1
    if err!=nil {
        return -1
    }      
    return res
}

///////////FUNCIONES DE LAS OPCIONES DEL MENÚ
//////////////////////

///OPCIÓN LISTAR
// Esta función hace la petición a la api y muestra las fotos
func ListarImagenes()(exito string,err error) {

    images, err := GetList()
     
    if(err==nil){

        for _, image:= range images {
            
            imageString := fmt.Sprintf("%d  -  %s", image.Id, image.Title);
            fmt.Println(imageString)
        }

        exito = "Operación exitosa"
    }

    return ;

}

///OPCIÓN CREAR
//Esta función pide el título y url de la imagen. Luego hace un post a la api para crear la imagen
func CrearImagen()(exito string, err error) {
    var title string
    var url string

    fmt.Println("Escribe el título: ")
    fmt.Scanln(&title)

    fmt.Println("Escribe la url: ")
    fmt.Scanln(&url)

    image := Instanciarimagen(0, title, url)

    
    err = PostImage(image)

    if(err==nil){
        exito = "Imagen creada"
    }

    return 
}

// Esta función crea y llena la estructura de la Imagen
func Instanciarimagen(id int, title string, url string) Image{

    image := Image{}
    image.Title=title
    image.AlbumId= 1
    image.Url = url
    image.ThumbnailUrl = url

    return image
}

//OPCIÓN MODIFICAR
//Esta función pide el id de la imagen a modificar y hace una petición PUT en la api
func ModificarImagen()(exito string, err error){
    var title string
    var url string
    var id int


    fmt.Println("Escribe el id: ")
    fmt.Scanln(&id)

    fmt.Println("Escribe el título: ")
    fmt.Scanln(&title)

    fmt.Println("Escribe la url: ")
    fmt.Scanln(&url)

    image := Instanciarimagen(id, title, url)

    
    err = PutImage(image)

    if(err==nil){
        exito = "Imagen modificada"
    }

    return 
}

//OPCIÓN CONSULTAR()
//Esta petición pide la id de la imagen y hace un get a la api usando ese id
func ConsultarImagen() (exito string, err  error){
    var id int

    fmt.Println("Escribe el id: ")
    fmt.Scanln(&id)

    imag, err := GetOne(id)


    if(err== nil){
        fmt.Println("Id: ", imag.Id)
        fmt.Println("Título: ", imag.Title)
        fmt.Println("URL: ", imag.Url)
        fmt.Println("URL (chica): ", imag.ThumbnailUrl)

        exito = "Operación exitosa"
    }

    return 
}

//OPCIÓN BORRAR
//Esta función pide un id y hace una petición delete a la api
func BorrarImagen() (exito string,err error){
    var id int

    fmt.Println("Escribe el id: ")
    fmt.Scanln(&id)

    err = DeleteImage(id)

    if(err == nil){
        exito = "Imagen borrada"
    }

    return ;
}

/////////// FUNCIONES QUE INTERACTÚAN CON LA API
//////////////////////

//Esta función hace un get a la api y retorna un arreglo de estructuras con las imagenes
func GetList() ([]Image, error) {

    res, err := Get(API_URL, 0);
    var images []Image

    if(err==nil){
        err = json.Unmarshal(res,&images)
    }
   
    return images, err
}

//Esta función recibe una imagen y hace un post hacía la api
func PostImage(image Image) error{
    return Post(API_URL, image);
}

//Esta función hace un get usando el id que recibe como parámetro
func GetOne(id int) (Image, error){
    res, err := Get(API_URL, id);
    var image Image

    err = json.Unmarshal(res,&image)

    return image, err
    
}

//Esta función hace un get sobre la url indicada, si el id es mayor a cero, hace el get incluyendo el id
func Get(url string, id int) ([]byte, error){

    fullUrl := url

    if id > 0{
        fullUrl = fmt.Sprintf("%s/%d",fullUrl, id)
    }

    response, err := http.Get(fullUrl)
    
    if(err ==nil){    
        responseData, err := ioutil.ReadAll(response.Body)
        ImprimirRespuesta(response.Body)
        return responseData, err
    }

    return nil, err
}


//Esta función hace un post sobre la url y la imagen indicada
func Post(url string, image Image)(error) {
    jsonData, err := json.Marshal(image)
    res, err := http.Post(url,"application/json", bytes.NewBuffer(jsonData))

    if(err==nil){
        ImprimirRespuesta(res.Body)
        PostSlack("Imagen creada, título: " + image.Title)
        
    }
   
    return err
}

//Esta función hace una petición delete usando el id indicado
func DeleteImage(id int) error{

    url := fmt.Sprintf("%s/%d", API_URL, id );
    res, err := http.NewRequest("DELETE", url, nil);

    if(err==nil){
        ImprimirRespuesta(res.Body)
        PostSlack("Imagen eliminada Id: " + strconv.Itoa(id))
    }

    return err
}

//Esta función hace una petición put usando la imagen enviada
func PutImage(image Image) error{

    url := fmt.Sprintf("%s/%d", API_URL, image.Id);
    jsonData, err := json.Marshal(image)
    res, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))

    if(err==nil){
        ImprimirRespuesta(res.Body)
    }

    return err
}

//Esta función manda un mensaje a slack
func PostSlack(mensaje string){

    json :=  []byte(`{"text" : "` + mensaje + `"}`)
    _, err := http.Post(API_SLACk,"application/json", bytes.NewBuffer(json))

    if err!=nil{
        fmt.Println("Error al enviar a slack")
        fmt.Println(err)
    }
}

func ImprimirRespuesta(res io.ReadCloser){
    if res==nil {return }

    responseData, _ := ioutil.ReadAll(res)
    
    
    mensaje := string(responseData)
    
    if len(mensaje) > 0{
        fmt.Print("\n\nRespuesta del servidor\n")
        fmt.Print(mensaje+ "\n\n")
    }
}



package main


import (
    "strconv"
    "bytes"
    "encoding/json"
	"fmt"
	"io/ioutil"
    "net/http"
)


//Constantes
const API_URL = "https://jsonplaceholder.typicode.com/photos"

//Opciones
const (
    OpcionListar = 1
    OpcionCrear = 2
    OpcionModificar  = 3
    OpcionConsultar =4
    OpcionEliminar = 5
    OpcionSalir = 6
)






type Image struct{
    AlbumId int  `json:albumId` 
    Id int  `json:id` 
    Title string  `json:title` 
    Url string  `json:url` 
    ThumbnailUrl string  `json:thumbnailUrl` 
}



func main() {



    opcion := 1


    for {
        opcion = Menu()

        var error error

        if opcion == OpcionListar {
            error = ListarImagenes()
            continue
        }



        if opcion == OpcionCrear {
            error = CrearImagen()
        }



        if opcion == OpcionModificar {
            error = ModificarImagen()
        }



        if opcion == OpcionConsultar {
            error = ConsultarImagen()
        }


        if opcion == OpcionEliminar {
            error = BorrarImagen()
        }

        if opcion == OpcionSalir {
            break
        }

        if error != nil {
            fmt.Println(error)
            continue
        }

        if !esOpcionValida(opcion){
            fmt.Println("Opción no válida")
        }
    }

}

func esOpcionValida(opcion int) bool{
    return opcion>0 && opcion<7
}


func Menu() int {
    fmt.Println("Elige una opción:")

    fmt.Println("1. Listar imagenes")
    fmt.Println("2. Crear imagen")
    fmt.Println("3. Modificar imagen")
    fmt.Println("4. Consultar una imagen:")
    fmt.Println("5. Eliminar imagen")
    fmt.Println("6. Salir")

    var input string
    fmt.Scanln(&input)

    res, err := strconv.Atoi(input);

    if err!=nil {
        return -1
    }
        

    return res

}




///OPCIÓN LISTAR
func ListarImagenes() error {

    imagenes, error := GetList()

    fmt.Println(imagenes)
    
    if(error!=nil){

       
        for _, imagen := range imagenes {
            imageString := fmt.Sprintf("%d  -  %s", imagen.Id, imagen.Title);
            fmt.Println()
        }
    }

    return error;

}




///OPCIÓN CREAR
func CrearImagen()error {
    var title string
    var url string

    fmt.Println("Escribe el título: ")
    fmt.Scanln(&title)

    fmt.Println("Escribe la url: ")
    fmt.Scanln(&url)

    image := Instanciarimagen(0, title, url)

    
    err:= PostImage(image)

    return err
}




func Instanciarimagen(id int, title string, url string) Image{

    image := Image{}
    image.Title=title
    image.AlbumId= 1
    image.Url = url
    image.ThumbnailUrl = url

    return image
}





//OPCIÓN MODIFICAR
func ModificarImagen() error{
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

    
    err:= PutImage(image)

    return err
}


//OPCIÓN CONSULTAR()
func ConsultarImagen() error{
    var id int

    fmt.Println("Escribe el id: ")
    fmt.Scanln(&id)

    imag, err := GetOne(id)


    if(err== nil){
        fmt.Println("Id: %d", imag.Id)
        fmt.Println("Título: %s", imag.Title)
        fmt.Println("URL: %s", imag.Url)
        fmt.Println("URL (chica): %s", imag.ThumbnailUrl)
    }



    return err
}

//OPCIÓN BORRAR
func BorrarImagen() error{
    var id int

    fmt.Println("Escribe el id: ")
    fmt.Scanln(&id)

    err := DeleteImage(id)

    return err;

}

func GetList() ([]Image, error) {

    res, err := Get(API_URL, 0);
    var images []Image

    if(err==nil){
        err = json.Unmarshal(res,&images)
    }
   

    return images, err
}

func PostImage(image Image) error{
    return Post(API_URL, image);
}


func GetOne(id int) (Image, error){
    res, err := Get(API_URL, id);
    var image Image

    err = json.Unmarshal(res,&image)

    return image, err
    
}












/***
PRINCIPALES
*/

func Get(url string, id int) ([]byte, error){

    fullUrl := url

    if id > 0{
        fullUrl = fmt.Sprintf("%s/%d",fullUrl, id)
    }

    response, err := http.Get(fullUrl)
    
    if(err ==nil){    
        responseData, err := ioutil.ReadAll(response.Body)
        return responseData, err
    }

    return nil, err
}








func Post(url string, image Image)(error) {

    jsonData, err := json.Marshal(image)
    _, err = http.Post(url,"application/json", bytes.NewBuffer(jsonData))


    return err
}


func DeleteImage(id int) error{

    url := fmt.Sprintf("%s/%d", API_URL, id );
    _, err := http.NewRequest("DELETE", url, nil);

    return err
}


func PutImage(image Image) error{

    url := fmt.Sprintf("%s/%d", API_URL, image.Id);
    jsonData, err := json.Marshal(image)
    _, err = http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))

    return err
}






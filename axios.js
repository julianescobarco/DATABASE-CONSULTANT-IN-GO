const { default: axios } = require("axios");

/*axios.get(`https://api.themoviedb.org/3/movie/popular?api_key=192e0b9821564f26f52949758ea3c473&language=es-MX`)
.then((respuesta) => {
    //console.log('Funciona')
    console.log(respuesta.data.results)

})
.catch((error) => {
    console.log('No funciona hay algun error! '+ error)
})*/

const obtenerPelic = async() => {
    try{
        const respuesta = await axios.get('https://api.themoviedb.org/3/movie/popular?', {
            params: {
                api_key:'192e0b9821564f26f52949758ea3c473',
                language: 'es-MX'
            }
        })
        console.log(respusta);

    }catch(error){
        console.log(error);

    }
}

obtenerPelic();
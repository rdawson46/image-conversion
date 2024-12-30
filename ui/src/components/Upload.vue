<script setup>
    import {ref} from 'vue';

    function makeParagraph(line) {

    }


    function uploadFile() {
        const imagePreview = document.getElementById('imagePreview');
        const imageResult = document.getElementById('imageResult');
        const input = document.getElementById('image');
        let files = input.files;

        if (files.length) {
            const reader = new FileReader();
            let file = files[0];

            reader.onload = (e) => {
                imagePreview.src = e.target.result
                const myForm = document.getElementById('myForm');
                const formData = new FormData(myForm)

                fetch('/upload', {
                    method: 'POST',
                    body: formData
                })
                .then(response => {
                    if (!response.ok) {
                        throw new Error("Server response broke")
                    }
                    // return response.blob()
                    return response.text()
                })
                .then(text => {
                    /*
                    need to split text on new line and create paragraph per line
                    */
                    const lines = text.split('\n')

                    imageResult.textContent = text
                })
                /*
                .then(blob => {
                    // temp
                    console.log("state")
                    const url = window.URL.createObjectURL(blob);

                    const link = document.createElement('a');
                    link.href = url;
                    link.download = 'testing.gzip';

                    document.body.appendChild(link);
                    link.click();
                    document.body.appendChild(link);
                })
                */
                .catch(error => {
                    console.log(error)
                });
            }

            reader.readAsDataURL(file);
        }
    }
</script>

<template>
    <form id="myForm" v-on:submit.prevent="uploadFile">
        <label for="upload">Upload image</label>
        <input id="image" name="image" type="file" accept="image/*"></input>

        <button>Submit</button>
    </form>

    <img id="imagePreview"></img>
    <div id="imageResult"></div>
</template>

<style scoped>
img {
    padding: 16px;
    width: 30rem;
}
</style>

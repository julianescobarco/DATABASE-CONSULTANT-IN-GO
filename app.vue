<template>
    <div>
        <p> Texto component</p>
        <div>{{ text }}</div>
    </div>
</template>

<script>
    export default {
        data() {
            return{ 
                text: "Hola!!.."
            }
        }
    };
</script>

<style>
p {
    color:blue;
}

/** Next course Composition API */
</style>
<?php
    //get random image
    $files = glob('images/*.{jpg,jpeg,png,gif}', GLOB_BRACE);
    $img = $files[array_rand($files)];
    $img = 'images/' . basename($img);

    //send to client
    header('Content-Type: image/jpeg');
    header('Image-Source: ' . str($img));
    header('Server: Crime.cx Cat Generator');
    header('Content-Length: ' . filesize($img));
    readfile($img);
?>
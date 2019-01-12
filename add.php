<?php
    /*
     * @file:    add.php
     * @brief:   Retreives a new word from the formular of the mainpage,
     *           filters and hands it over to the Vocabulary-API
     * @author:  Raphael Pour <info@raphaelpour.de>
     * @date:    01-2019
     */

    require_once 'Vocabulary.php';
    $de = filter_input(INPUT_POST, 'de', FILTER_SANITIZE_SPECIAL_CHARS);
    $en = filter_input(INPUT_POST, 'en', FILTER_SANITIZE_SPECIAL_CHARS);

    if($de && $en)
    {   
        $voc = new Vocabulary();
        $voc->addWord($de, $en);
    }
    include 'index.php'; 
?>

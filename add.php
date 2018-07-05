<?php
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

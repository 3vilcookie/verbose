<?php
require_once 'Vocabulary.php';

$voc = new Vocabulary();

$randomTranslation = $voc->getRandomWord();
$randomWord = $randomTranslation['en'] . " - " . $randomTranslation['de'];
?>
<!DOCTYPE html>
<html lang="en">
<head>
<title>Verbose</title>
<meta charset="utf-8">
<link rel="shortcut icon" type="image/png" href="logo.png">
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/css/bootstrap.min.css" integrity="sha384-WskhaSGFgHYWDcbwN70/dfYBj47jz9qbsMId/iRN3ewGhXQFZCSftd1LZCfmhktB" crossorigin="anonymous">
</head>
<body>
<div class='container jumbotron' style='padding-top:30px;padding-bottom:30px;'>
<div class="row">
    <div class="col-lg-2"><a href='.'><img src='logo.png'></a></div>
    <div class="col-lg-4">
        <h1><em>verbose</em></h1>
        <p>wortreich, langatmig, ausführlich, weitschweifig</p>
    </div>
</div>
</div>

<div class="container">
<div class="alert alert-success">Word of the Pageload:  <strong><?php echo $randomWord; ?></strong></div>
</div>

<div class="container">
    <form method="post" action="./add.php">
     <div class="form-group row">
        <div class="col-lg-2">
            <input class="form-control" id="en" name="en" placeholder="verbose" type="text"/>
        </div>
        <div class="col-lg-2">
            <input class="form-control" id="name1" name="de" placeholder="quasseln" type="text"/>
        </div>
        <div class="col-lg-2">
                <button class="btn btn-success " name="submit" type="submit">Add</button>
        </div>
    </div>
    </form>
   </div>
</div>
<div class='container'>
<table class='table table-striped'>
<thead>
    <tr>
    <th>EN</th>
    <th>DE</th>
    </tr>
</thead>
<tbody>
<?php
$words = $voc->getWordList();
$count = $voc->getWordCount();

foreach($words as $word)
    echo "<tr><td>" . $word['en'] . "</td><td>" . $word['de'] . "</td></tr>\n";   
?>
</tbody>
</table>
<div>Words: <?php echo $count; ?></div>
</div>

<footer>
<div class='container' style='text-align:center'>
<hr>
(C) 2018 <a href='https://raphaelpour.de'>Raphael Pour</a>| <a href='https://www.gnu.org/licenses/lgpl-3.0.en.html'>LGPL</a> | <a href='https://raphaelpour.de/impressum/'>Impressum</a>
</div>
</footer>

<script>
document.getElementById("en").focus();
</script>

</body>
</html>

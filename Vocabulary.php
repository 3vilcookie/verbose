<?php
    /*
     * @file:    Vocabulary.php
     * @brief:   Vocabulary-API with db access
     * @author:  Raphael Pour <info@raphaelpour.de>
     * @date:    01-2019 - 07-2020
     */

ini_set("log_errors", 1);
ini_set("error_log","/var/log/verbose_error.log");

class Vocabulary
{

    private $user = "";
    private $password = "";

    private $db = Null;

    public function __construct()
    {
        require "Config.php";
        $this->password = $CONFIG['db_password'];
        $this->user = $CONFIG['db_user'];
        $this->db = new PDO('mysql:host=localhost;dbname=langpp;charset=utf8', $this->user, $this->password);
    }

    public function getWordCount()
    {
        $statement = $this->db->prepare("SELECT COUNT(*) as count FROM vocabulary");

        if($statement->execute())
          return $statement->fetch()['count'];

        error_log("Error getting the word count: " . $statement->errorInfo());        
        return 0;
    }

    public function getSortedWordList($column, $cending = true)
    {
        $statement = $this->db->prepare("SELECT * FROM vocabulary ORDER BY :column :cending");
        $cendingStr = ($cending) ? ("ASC") : ("DSC");
    
        $statement->bindParam(':column', $column);
        $statement->bindParam(':cending', $cendingStr);

        $list = array();

        if($statement->execute())
            while($row = $statement->fetch())
                array_push($list, array('en' => $row['en'], 'de' => $row['de']));

        error_log("Error getting sorted word list: " . $statement->errorInfo());        
        return $list;
    }

    public function getWordList()
    {
        $statement = $this->db->prepare("SELECT * FROM vocabulary ORDER BY en");

        $list = array();

        if($statement->execute())
            while($row = $statement->fetch())
                array_push($list, array('en' => $row['en'], 'de' => $row['de']));

        error_log("Error getting word list: " . $statement->errorInfo());        
        return $list;
    }

    public function addWord($de, $en)
    {
        $statement = $this->db->prepare("INSERT INTO vocabulary (de, en) VALUES(:de,:en)");

        $statement->bindParam(':de', $de,PDO::PARAM_STR );
        $statement->bindParam(':en', $en,PDO::PARAM_STR);

        if($statement->execute() !== True)
            error_log("Error adding a word " . $statement->errorInfo());        
    }

    public function getRandomWord()
    {   
        $statement = $this->db->prepare("SELECT * FROM vocabulary ORDER BY RAND() LIMIT 1");

        if($statement->execute())
        {
            $row = $statement->fetch();
            return array('en' => $row['en'], 'de' => $row['de']);
        }
        
        error_log("Error getting a random word " . $statement->errorInfo());        
        return array();
    }
}
?>

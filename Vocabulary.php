<?php


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
        else
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

        return $list;
    }

    public function getWordList()
    {
        $statement = $this->db->prepare("SELECT * FROM vocabulary ORDER BY en");

        $list = array();

        if($statement->execute())
            while($row = $statement->fetch())
                array_push($list, array('en' => $row['en'], 'de' => $row['de']));

        return $list;
    }

    public function addWord($de, $en)
    {
        $statement = $this->db->prepare("INSERT INTO vocabulary (de, en) VALUES(:de,:en)");

        $statement->bindParam(':de', $de,PDO::PARAM_STR );
        $statement->bindParam(':en', $en,PDO::PARAM_STR);

        if($statement->execute() !== True)
        {
            var_dump( $statement->errorInfo() );
        }
    }

    public function getRandomWord()
    {   
        $statement = $this->db->prepare("SELECT * FROM vocabulary ORDER BY RAND() LIMIT 1");

        if($statement->execute())
        {
            $row = $statement->fetch();
            return array('en' => $row['en'], 'de' => $row['de']);
        }
        else
        {
            return array('en' => 'DATABASE', 'de' => 'FUCKUP');
        }


    }

}
?>

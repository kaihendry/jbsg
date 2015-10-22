<?php
// Converts XML to JSON. Nightmare.

date_default_timezone_set('Asia/Singapore');
header('Content-Type: application/json');

// Keys suck http://dabase.com/blog/Javascript_API_barriers/
$creds = parse_ini_file(".creds.ini");
$xml = simplexml_load_file("http://www.nea.gov.sg/api/WebAPI/?dataset=psi_update&keyref=" . $creds["key"]);

// When debugging save the XML to /tmp/xml.txt
//$xml = simplexml_load_file("/tmp/xml.txt");

$array = array();

foreach ($xml->item->region as $region) {
	if ($region->id == "rNO") {
		//print_r($region->record);
		$attr = $region->record->attributes();
		foreach($attr as $key=>$val){
			$array[(string)$key] = date("c", strtotime ((string)$val));
		}
		foreach ($region->record->reading as $reading) {
			$attr = $reading->attributes();
			$t = (string) $attr->type;
			$v = (string) $attr->value;
			$array[$t] = $v;
		}
		break;
	}
}

echo json_encode($array, JSON_PRETTY_PRINT);

?>

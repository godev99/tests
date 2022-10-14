****BADGE READER - EXERCICE TECHNIQUE****

La solution proposée au test de 'lecteur de badge' a été réalisée en Go.
Comme suggéré, un module a été développé. Celui sera consommé par un client sous forme de dépendance.
Ce test comprends également des tests unitaires pour l'ensemble des fonctions du module.


## Table of contents
1. [CLIENT](#CLIENT)
2. [MODULE](#MODULE)
3. [TESTS UNITAIRES](#TestsUnitaires)
4. [EXECUTION](#Execution)

---
### CLIENT <a name="CLIENT"></a>


Import du module dépendant.

```
import (
	...
	"greengo/datachunk"
	...
)
```

Une variable globale est nécessaire :

- chunks : Contient l'ensemble des jeux de données, construit au fil des appels au module.


Plutôt que d'avoir une saisie utilisateur, ces jeux de données ont été préalablement intégrés au client :

```
	// chunks of data
	slice := [][]int{[]int{10, 20, 30, 40, 65}, []int{65, 66}, []int{67, 100, 200, 30, 40, 50}, []int{88, 65}, []int{65, 66, 36, 3, 95, 5}}

```

Les bornes de recherche sont également déclarées comme client, en tant que chaines de caractères. Le module les convertira en chaine d'entiers.

```
	// begin and end sequence to identify badge
	startCharacters := []int{"SOH", "STX"} // ABC
	endCharacters   := []int{"ETX", "EOT"}   // XYZ
```

Par la suite, nous opérons une boucle sur les jeux de données, en passant en paramètres :

- Le jeu de données courant.
- Le tableau d'entier coté client. Celui constitue le lien entre chaque appel, car il concatène l'ensemble des jeux de données.
- les délimiteurs (afin d'éviter de les intégrer coté module).

```
	// here we will call GetDataChunk function from our own private module greengo/datachunk
	for i := range slice {
		badge, isBadgeFound := datachunk.GetDataChunk(slice[i], &chunks, startCharacters, endCharacters)
		if isBadgeFound {
			fmt.Println("chunks:", chunks)
			fmt.Println("badge found:", badge)
			os.Exit(1)
		}
	}
	fmt.Println("badge not found")
	fmt.Println("chunks:", chunks)
```

Dès qu'un badge est identifié, le client stop les appels et la variable contenant l'ensemble de jeux de données est réinitialisée.
Sinon, l'ensemble du jeu de données est affiché et ensuite réinitialisé.


---
### MODULE <a name="MODULE"></a>

La méthode exportée "GetDataChunk" permet de traiter des jeux de données, un par un, en incrémentant une variable globale coté client.
À chaque envoi de données, un appel est fait à une fonction privée (getSubString) qui se charge de vérifier la présence d'un badge.

***Ce type implémentation ne convient pas à un environnement de production. L'idéal étant de scanner le jeu de données uniquement si x valeurs sont détectées.
Ceci est plus couteux en termes de code. Pour les besoins de l'exercice, l'implémentation ne prend pas en compte de telles contraintes.***

Finalement, 2 valeurs sont renvoyées au client :

- myBadge : est soit vide, soit contient le contenu d'un badge
- isMyBadgeFound : A-t-on trouvé un badge ?

### GetDataChunk

Voici les opérations effectuées chronologiquement par cette fonction :
- Concaténer le nouveau jeu de données avec les précédents.
- Le convertir en chaine de caractères.
- Convertir les bornes en chaine d'entiers.
- Analyser et récupérer un badge qui existerait dans le jeu de données concaténées.
- Convertir la chaine de caractères en tableau d'entiers.
- Renvoyer le statut de la recherche de badge ainsi que ce dernier, qu'il ait une valeur ou non.

_Plutôt que de traiter des recherches sur des tableaux d'entiers (pas de fonction native en go), l'approche de convertir en string permet de simplifier le traitement et de rendre le code plus maintenable_

```
// GetDataChunk : Main function will provide a badge if exist with chunks
func GetDataChunk(chunk []int, chunks *[]int, startWith []string, endWith []string) (myBadge []int, isMyBadgeFound bool) {

	// variable for badge in string and bounds conversion from ascii to string
	var badgeInString, downBoundary, upBoundary string
	// let's add current chunk to all chunks
	for k := range chunk {
		*chunks = append(*chunks, chunk[k])
	}

	// transform list of int into string to facilitate badge detection
	ChunksInString := sliceToString(*chunks)

	// Convert boundaries ascii into two strings
	for i := range startWith {
		number := startWith[i]
		downBoundary += AsciiToBase10(number) + " "
	}
	downBoundary = strings.TrimRight(downBoundary, " ")

	for i := range endWith {
		number := endWith[i]
		upBoundary += AsciiToBase10(number) + " "
	}
	upBoundary = strings.TrimRight(upBoundary, " ")

	// If badge is not found, badge variable is empty
	badgeInString, isMyBadgeFound = getSubString(ChunksInString, downBoundary, upBoundary)

	// revert badge string to slice of int
	myBadge = stringToSlice(badgeInString)

	return myBadge, isMyBadgeFound
}
```
### sliceToString
Cette fonction récupère une chaine de caractères à partir d'un tableau d'entiers à partir d'une chaine de caractères.

```
func sliceToString(mySlice []int) (myList string) {

	// this will host final string
	var valuesText []string

	for i := range mySlice {
		number := mySlice[i]
		text := strconv.Itoa(int(number))
		valuesText = append(valuesText, text)
	}

	result := strings.Join(valuesText, " ")

	// We need to separate each new chunk with a space as soon as we concat multiple chunks
	if myList == "" {
		myList = result
	} else {
		myList += " " + result
	}

	return myList
}
```

### stringToSlice
Cette fonction récupère un tableau d'entiers à partir d'une chaine de caractères.

```
func stringToSlice(myString string) (mySlice []int) {

	// this will host slice result
	numbers := strings.Fields(myString)

	for i := range numbers {
		number, _ := strconv.Atoi(numbers[i])
		mySlice = append(mySlice, number)
	}

	return mySlice
}
```

### getSubString
Avec les différents paramètres convertis en chaines de caractères, nous allons traiter le jeu de caractères global comme une chaine :

- Recherche de l'index de la borne inférieure dans le jeu de données global.
- On initialise un nouveau tableau en partant de l'index trouvé jusqu'à la fin du jeu de données.
- En cas d'index négatif, on sort avec le statut found = false et la valeur de myBadge à vide.
- Recherche de l'index de la borne supérieure à partir du nouveau tableau.
- En cas d'index négatif, on sort avec le statut found = false et la valeur de myBadge à vide.
- On redécoupe le nouveau tableau à partir du début jusqu'au nouvel index trouvé.
- Si le tableau ne contient aucun élément, on sort avec le statut found = false et la valeur de myBadge à vide.
- On renvoie le badge trouvé ainsi que le statut found à true.


```
func getSubString(myString string, startBound string, endBound string) (myBadge string, found bool) {

	// Searching for occurrence of first upBoundary
	sliceWithStartBound := strings.Index(myString, startBound)
	// startBound was not found, exit
	if sliceWithStartBound == -1 {
		return "", false
	}

	// New string from startBound to end of original string, myString
	newString := myString[sliceWithStartBound+len(startBound):]

	sliceWithStartandEndBounds := strings.Index(newString, endBound)
	if sliceWithStartandEndBounds == -1 {
		return "", false
	}

	// to dig why we should begin slice from 1 to avoid empty slice with len = 1 for example (thanks to unit test)
	myBadge = newString[1:sliceWithStartandEndBounds]

	// if two bounds were founds but no element in between, this should be trapped (thanks to unit test)
	if len(myBadge) == 0 {
		return "", false
	}

	return myBadge, true
}
```

### TestAsciiToBase10

_Les expressions régulières auraient résolu le problème, mais le compromis entre l'efficacité et la lisibilité parait moins avantageux._

---
### TESTS UNITAIRES <a name="TestsUnitaires"></a>

Ces tests couvrent les 4 fonctions du module :

- TestGetDataChunk > fonction GetDataChunk
- TestSliceToString > fonction sliceToString
- TestStringToSlice > fonction sliceToString
- TestGetSubString > fonction getSubString

Pour chacun des tests, plusieurs jeux de paramètres ont été donnés.

### TestGetDataChunk
#### Paramètres
- name: le nom du test (statut prévu)
- startWith: la borne inférieure
- endWith: la borne supérieure
- chunk: le jeu de données

_Le jeu de données ne fait pas partie de la structure de données, car il est initialisé au début du test._

```
{ name: "Ok", startWith: []int{65}, endWith: []int{66}, chunk: []int{65, 20, 30, 40, 66} }
{ name: "Ok", startWith: []int{65}, endWith: []int{66}, chunk: []int{10, 20, 65, 66, 65} }
{ name: "Ko", startWith: []int{65, 45}, endWith: []int{66}, chunk: []int{10, 66, 65, 40, 65} }
{ name: "Ok", startWith: []int{65}, endWith: []int{66}, chunk: []int{10, 6, 65, 40, 6} }
{ name: "Ko", startWith: []int{65, 57}, endWith: []int{66, 95}, chunk: []int{65, 10, 57, 20, 30, 40, 65} }
```
#### Résultats attendus
La présence d'un badge dans le jeu de données.

### TestSliceToString
#### Paramètres
- name: nom du test (statut prévu)
- mySlice : le tableau d'entiers

```
{ name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 20, 30, 40, 3} }
{ name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 1, 65, 66, 65} }
{ name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{1, 66, 65, 40, 2} }
{ name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 6, 65, 2, 6} }
{ name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 10, 1, 3, 2, 4, 65} }
{ name: "KO", startWith: []string{"SOH"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 3, 30, 40, 1} }
{ name: "KO", startWith: []string{"SOH"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 20, 65, 3, 65} }
{ name: "OK", startWith: []string{"SOH"}, endWith: []string{"ETX", "EOT"}, chunk: []int{1, 66, 65, 3, 4} }
{ name: "KO", startWith: []string{"SOH"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 4, 65, 1, 6} }
{ name: "OK", startWith: []string{"STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 10, 1, 20, 30, 3, 4} }
{ name: "KO", startWith: []string{"STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 1, 30, 40, 66} }
{ name: "KO", startWith: []string{"STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{1, 20, 65, 66, 65} }
{ name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 66, 65, 40, 65} }
{ name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 6, 65, 40, 6} }
{ name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 10, 57, 20, 30, 40, 65} }

```

#### Résultats attendus
Une chaine de caractères qui n'est pas vide

### TestStringToSlice
#### Paramètres
- name: nom du test (statut prévu)
- myString : la chaine de caractères

```
{ name: "65, 54, 678, 23, 0, 1", mySlice: []int{65, 54, 678, 23, 0, 1} }
{ name: "-1, 9, 54, 0, 54", mySlice: []int{-1, 9, 54, 0, 54} }
{ name: "0", mySlice: []int{0} }
{ name: "235334, 564, 6776", mySlice: []int{235334, 564, 6776} }
	
```

#### Résultats attendus
Une chaine de caractères qui n'est pas vide

### TestGetSubString
#### Paramètres
- name: nom du test (statut prévu)
- myString : la chaine de caractères
- startWith: la borne inférieure
- endWith: la borne supérieure

```
{ name: "65 54 678 23 0 1 (65 1)", myString: "65 54 678 23 0 1", startWith: "65", endWith: "1" }
{ name: "1 9 54 0 54 (9 -1)", myString: "-1 9 54 0 54", startWith: "9", endWith: "-1" }
{ name: "0 (65 1)", myString: "0", startWith: "65", endWith: "1" }
{ name: "12 434 65 564 6776 (65 564)", myString: "12 434 65 564 6776", startWith: "65", endWith: "564" }
{ name: "43 65 56 67 45 12 (56 12)", myString: "43 65 56 67 45 12", startWith: "56", endWith: "12" }
{ name: "43 65 56 67 45 12 (56 45)", myString: "43 65 56 67 45 12", startWith: "56", endWith: "45" }

```

#### Résultats attendus
Une sous chaine de caractères qui n'est pas vide

### TestAsciiToBase10
#### Paramètres
- name: nom du test (statut prévu)
- myString : la chaine de caractères

```
{ myString: "SOH" }
{ myString: "STX" }
{ myString: "EOT" }
{ myString: "EQ" }
{ myString: "?" }
{ myString: "ETX" }
{ myString: "NUL" }
{ myString: "ENQ" }
{ myString: "LF" }
{ myString: "LIFO" }
{ myString: "'" }
{ myString: "DLE" }
{ myString: "DC1" }
{ myString: "SYN" }
{ myString: "LIO" }
{ myString: "}" }
{ myString: "EM" }
{ myString: "NAK" }
{ myString: "CAN" }
{ myString: "RLT" }
{ myString: "ESC" }
```

#### Résultats attendus
Un code ascii valide. Sachant qu'un code supérieur à 33, et dont la chaine de caractère est supérieure à 2, ne correspond pas aux caractères ascii non-imprimables.

---
### EXECUTION <a name="Execution"></a>
```
go run main.go
```
Pour les tests unitaires :
```
cd datachunk
go test -v
```

package main 
/*


change front and back naming convention..



*/
import (
	 "fmt"
	 // "io"
	 "math"
	 // "bufio"
	 "os"
	 "strconv"
	 "math/rand"

)

// Neuron data type
type Neuron struct{
	id int 						// these only here to track process
	weight float64
	nonActivatedWeight float64
	connections [] *Synapse
	
}
// Synapses are the connections between neurons
type Synapse struct{
	front * Neuron
	back * Neuron
	value float64
}

func main() {

	var input = os.Args[1:]
	fmt.Println(input);
	/*  arguments : 
		aka the weights of input nodes..
		first value --		
		second value --
	
	*/
	// going to assume perfect output right now --
	
	// check if the right amount of arguments were inputted..? "inputed"?
	fmt.Println("Ayo - Creating Neural Network");
	var in ,hid ,out = initNetwork(input);
	// MAKE IT SO YOU CAN CHANGE ITERATIONS DEPENDING ON INPUT******TODO
	 for i:=0; i<5;i++{
		for _ , element := range hid{
			for index , ele := range element.connections{
				fmt.Printf("%d : %f \n", index, ele.value);
			}
			//printSynapseValues(element.connections);
		}
		runSimulation(in,hid,out, input);
		// for _ , element := range hid{
		// 	printSynapseValues(element.connections);
		// }
	 }

	// reader := bufio.NewReader(os.Stdin)
 //    fmt.Print("Enter text: ")
 //    text, _ := reader.ReadString('\n')
 //    fmt.Println(text)
}

// intialise all elements needed for network to work
// may have to take input

func printSynapseValues(synapses [] Synapse){
	for index , ele := range synapses{
		fmt.Printf("%d : %f \n", index, ele.value);
	}
}
func initNetwork(inputValues [] string )([] * Neuron, [] * Neuron, [] * Neuron){
	// create input nodes -- 
	var hidNodesLen = 3
	var outNodesLen = 1
	// going to hardcore number of neurons for now ******
	hidNodes := make([]* Neuron, 0);
	inNodes := make([]* Neuron, 0);
	outNodes := make([]* Neuron, 0);
	// create input Neurons
	for index, element := range inputValues{
		var current = new(Neuron);
		current.id = index;
		weight, err := strconv.ParseFloat(element, 64);
		if(err != nil){
			fmt.Println(err);
		}
		current.weight = weight;
		current.nonActivatedWeight = 0;
		inNodes = append(inNodes, current);
	}
	// create hidden Neurons
	for  i:=0 ; i<hidNodesLen; i++ {
		var current = new(Neuron);
		current.id = i;
		current.weight = 0;
		current.nonActivatedWeight = 0;
		hidNodes = append(hidNodes, current);

	}
	// create output Neurons
	for i:=0; i<outNodesLen;i++{
		var current = new(Neuron);
		current.id = i;
		current.weight = 0;
		current.nonActivatedWeight = 0;
		outNodes = append(outNodes, current);
	}
	// add synpases?
	addSynapses(inNodes, hidNodes, outNodes);


	return inNodes, hidNodes, outNodes;
}
// dont really need a slice of out neurons if assuming there is only one rn.. but may be useful when scaling
func addSynapses(in [] *Neuron, hid [] *Neuron, out [] *Neuron){
	// var lenFrontSynapses = len(in)*len(hid);
	// adding the front synapses
	for i:=0; i<len(in); i++ {
		for j:=0; j<len(hid); j++{
			var synapse = new(Synapse);										// creating synapse/connection
			synapse.front = in[i];
			synapse.back = hid[j];
			synapse.value = myGaussian();
			in[i].connections = append(in[i].connections, synapse);		// adding synapse to both neurons
			hid[j].connections = append(hid[j].connections, synapse);
		}
	}
	// adding the back synapses
	// this assumes there is only one output rn..
	for i:=0; i<len(hid);i++{
		var synapse = new(Synapse);											// creating synapse/connection
		synapse.front = hid[i];
		synapse.back = out[0];
		synapse.value = myGaussian();
		hid[i].connections = append(hid[i].connections, synapse);			// adding synapse to both neurons
		out[0].connections = append(out[0].connections, synapse);
	}

}

func forwardProp(nodes [] *Neuron){
	for _ , node := range nodes{
		for _, connection := range node.connections{
			connection.back.weight += node.weight * connection.value;
		}
	}
}
func activateNodes(nodes [] *Neuron){
	for _ , node := range nodes{
		node.nonActivatedWeight = node.weight;
		// fmt.Printf("%d ---activated weight = %f\n", index, node.nonActivatedWeight);
		node.weight = sigmoidFunc(node.weight);
	}
	// printNonActivated(nodes);
}
// This backwardProp is for out Nodes to hidden nodes
func backwardPropOut(hidNodes [] *Neuron,outNodes [] *Neuron, deltaOutSum float64){
	// since this for a xor function..
	// input is a const
	// var expectedResult := strconv.ParseFloat(input[0], 64)^strconv.ParseFloat(input[1], 64);
	// fmt.Println("expectedResult = %f", expectedResult);
	// var calculated := outNodes[0]; 
	// var sumMarginOfError := expectedResult - calculated;

	// var deltaOutputSum := derivativeSigmoid(calculated) * sumMarginOfError;
	var deltaWeight = make([] float64,len(hidNodes));				// TODO!!!!! MAGIC
	//var outNode = outNodes[0];
	// finding the delta weights for the hidden nodes
	for index , node := range hidNodes{
		deltaWeight[index] =deltaOutSum * node.weight;
	}
	// adjusting new weights by adding delta weights to synapse values
	for index , connection := range outNodes[0].connections{
		connection.value = connection.value + deltaWeight[index];
	}



}
func backwardPropHid(inNodes [] *Neuron, hidNodes [] *Neuron, deltaHiddenSums [] float64, input [] string){
	var sizeOfdeltaWeights = len(input)*len(hidNodes);
	var hiddenDeltaWeights = make([]float64,sizeOfdeltaWeights) ;
	var counter =0;
	for _ , weight := range input{
		for index , _ := range deltaHiddenSums{ 
			var intermediate, _ = strconv.ParseFloat(weight, 64);// should really print out the err if its not nil..
			// fmt.Printf("hiddenDeltaWeights @ %d = %f * %f\n",index, intermediate, deltaHiddenSums[index]);
			hiddenDeltaWeights[counter] = intermediate*deltaHiddenSums[index];
			counter++;
		}
	}
	// fmt.Println("HiddenDelta Sums");
	// for index, weight := range deltaHiddenSums{
	// 	fmt.Printf("%d :: = %f\n", index, weight);
	// }

	// fmt.Println("HIDDEN DELTA WEIGHTS...")
	// for index , weight := range hiddenDeltaWeights{
	// 	fmt.Printf("%d :: = %f\n", index, weight);
	// }
	// adjusting new weight by adding delta weights to synapse values
	counter = 0; 		// reusing counter
	for _ , node := range inNodes {			// This could probably be a function frankly..
		for _ , connection := range node.connections{
			connection.value = hiddenDeltaWeights[counter] + connection.value;
		}
	}


}

func myGaussian() float64{
	// really handy function NormFloat64()
	var thingy = rand.NormFloat64();
	// fmt.Println("Number - %f" ,thingy );
	thingy = thingy * .5 + .5;
	// if the value is negative.. maybe get the absolute value not sure yet ^^
	return thingy;
}
// func gaussian() float64{
// 	fmt.Println("In gaussian");
//     var v, fac float64;
//     var phase = 0;
//     var S, Z, U1, U2, u float64;
//     S = 0;
//     if (phase==1){
//         Z = v * fac;
//     }else{
//        	for ok := true; ok; ok = (S >= 1){
//             U1 = float64(rand.Float64() / math.MaxFloat64);
//             U2 = float64(rand.Float64() / math.MaxFloat64);

//             u = 2. * U1 - 1.;
//             v = 2. * U2 - 1.;
//             S = u * u + v * v;
//         }

//         fac = math.Sqrt(-2. * math.Log(S) / S);
//         Z = u * fac;
//     }
//     phase = 1 - phase;
//     return math.Sqrt(Z*Z);
// }
// This is the activation function ;) Ill activate your function
func sigmoidFunc(x float64) float64{
    var intermediate = math.Exp(-x)+1;
    return 1/intermediate;
}

func derivativeSigmoid(x float64) float64{
	var intermediate = math.Pow((math.Exp(-x) + 1),2); 
	return math.Exp(-x)/ intermediate;
}
// Do this a thousand time.... and see how close we are.. then run it with different input while maintaining values
func runSimulation(in [] * Neuron, hid [] *Neuron, out [] *Neuron, input [] string){
	// run through process
	forwardSim(in,hid,out);
	backwardSim(in,hid,out, input);
	// forwardProp(in);
	// activateNodes(hid);
	// forwardProp(hid);
	// activateNodes(out);
	var first, _ = strconv.Atoi(input[0]);			// dont have any error handling right now...
	var second , _ = strconv.Atoi(input[1]);
	var expectedResult = float64(first^second);	// should print out err if not nil   **** MAY NOT WORK>>> as its an assertion..
	fmt.Printf("actual = %f\n", expectedResult);
	fmt.Printf("estimated == %f\n", out[0].weight);

}
//Simulates forward propagation
func forwardSim(in [] *Neuron, hid [] *Neuron, out [] *Neuron){
	forwardProp(in);
	// printNonActivated(hid);
	
	for _ , node := range hid{
		node.nonActivatedWeight = node.weight;
		// fmt.Printf("%d ---activated weight = %f\n", index, node.nonActivatedWeight);
		node.weight = sigmoidFunc(node.weight);
	}
	

	// activateNodes(hid);
	// printNonActivated(hid);

	forwardProp(hid);
	activateNodes(out);
}
// Simulates backWard propagation

// this should also probably take in or produce the delta output sum value
func backwardSim(in [] *Neuron, hid [] *Neuron, out [] *Neuron, input [] string){
	// gets delta output sum then uses it to calculate new values for synapses from output to hidden
	var first, _ = strconv.Atoi(input[0]);
	var second , _ = strconv.Atoi(input[1]);
	var expectedResult = float64(first^second);	// should print out err if not nil   **** MAY NOT WORK>>> as its an assertion..
	var calculated = out[0].weight; 
	var sumMarginOfError = expectedResult - calculated;
	var deltaOutputSum = derivativeSigmoid(calculated) * sumMarginOfError;
	
	// gets delta hidden sum then uses it to calculate new values for synapses from hidden to input
	// Delta hidden sum = delta output sum * hidden-to-outer weights * S'(hidden sum)

	// fmt.Println("Delta Output sum === %f", deltaOutputSum);
	// need to save synapses from output connections so they are not overwritten they will be needed later..
	copiedSynapsesValues := make([]*Synapse,len(out[0].connections) );
	copy(copiedSynapsesValues, out[0].connections);
	// for _ , element := range out{
	// 	printSynapseValues(element.connections);
	// }
	// fmt.Println("Second set of values ---\n");
	// for _ , elem := range copiedSynapsesValues{
	// 	fmt.Printf("%f ---\n",elem.value );
	// }
	// fmt.Println("		non activated values----\n");
	// for inde , elem := range hid{
	// 	fmt.Printf("%d ==== %f\n", inde , elem.nonActivatedWeight);
	// }
	
	backwardPropOut(hid,out, deltaOutputSum);
	
	var deltaHiddenSums = make([] float64,len(out[0].connections)) ;
	for index, _ := range deltaHiddenSums{
		// fmt.Printf("%f * %f * %f ---\n",deltaOutputSum, copiedSynapsesValues[index].value, hid[index].nonActivatedWeight);
		deltaHiddenSums[index]= deltaOutputSum * copiedSynapsesValues[index].value * hid[index].nonActivatedWeight;
	} 
	
	backwardPropHid(in, hid, deltaHiddenSums, input);

}

func printNonActivated(nodes [] *Neuron){
	fmt.Println("Printing non activated nodes\n");
	for index, elem := range nodes{
		fmt.Printf("%d ==== %f\n", index, elem.nonActivatedWeight);
	}
}

// func printNetwork(in [] *Neuron, hid [] *Neuron, out [] *Neuron){
// 	fmt.Println("-------------------Printing Network--------------------- ");
// 	for(){

// 	}
// }
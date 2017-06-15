package main 
/*


change front and back naming convention..



*/
import (
	 "fmt"
	 // "io"
	 "math"
	 "bufio"
	 "os"
	 "strconv"
	 "math/rand"

)

// Neuron data type
type Neuron struct{
	id int 						// these only here to track process
	weight float64
	nonActivatedWeight float64
	connections [] Synapse
	
}
// Synapses are the connections between neurons
type Synapse struct{
	front * Neuron
	back * Neuron
	value float64
}

func main() {

	const input := os.Args[1:]
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
	runSimulation(in,hid,out);


	reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter text: ")
    text, _ := reader.ReadString('\n')
    fmt.Println(text)
}

// intialise all elements needed for network to work
// may have to take input
func initNetwork(inputValues [] string )([] Neuron, [] Neuron, [] Neuron){
	// create input nodes -- 
	var hidNodesLen = 3
	var outNodesLen = 1
	// going to hardcore number of neurons for now ******
	hidNodes := make([] Neuron, 0);
	inNodes := make([] Neuron, 0);
	outNodes := make([] Neuron, 0);
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
		inNodes = append(inNodes, *current);
	}
	// create hidden Neurons
	for  i:=0 ; i<hidNodesLen; i++ {
		var current = new(Neuron);
		current.id = i;
		current.weight = 0;
		current.nonActivatedWeight = 0;
		hidNodes = append(hidNodes, *current);

	}
	// create output Neurons
	for i:=0; i<outNodesLen;i++{
		var current = new(Neuron);
		current.id = i;
		current.weight = 0;
		outNodes = append(outNodes, *current);
	}
	// add synpases?
	fmt.Println("About to addSynapses");
	addSynapses(inNodes, hidNodes, outNodes);


	return inNodes, hidNodes, outNodes;
}
// dont really need a slice of out neurons if assuming there is only one rn.. but may be useful when scaling
func addSynapses(in [] Neuron, hid [] Neuron, out [] Neuron){
	// var lenFrontSynapses = len(in)*len(hid);
	// adding the front synapses
	fmt.Println("length of in =%d", len(in));
	fmt.Println("length of hid =%d", len(hid));
	fmt.Println("length of out =%d", len(out));
	for i:=0; i<len(in); i++ {
		for j:=0; j<len(hid); j++{
			var synapse = new(Synapse);										// creating synapse/connection
			synapse.front = &in[i];
			synapse.back = &hid[j];
			synapse.value = myGaussian();
			fmt.Println("value : %f" ,synapse.value);
			in[i].connections = append(in[i].connections, *synapse);		// adding synapse to both neurons
			hid[j].connections = append(hid[j].connections, *synapse);
		}
	}
	// adding the back synapses
	// this assumes there is only one output rn..
	for i:=0; i<len(hid);i++{
		var synapse = new(Synapse);											// creating synapse/connection
		synapse.front = &hid[i];
		synapse.back = &out[0];
		synapse.value = myGaussian();
		fmt.Println("value :%f " , synapse.value);
		hid[i].connections = append(hid[i].connections, *synapse);			// adding synapse to both neurons
		out[0].connections = append(out[0].connections, *synapse);
	}

}

func forwardProp(nodes [] Neuron){
	for _ , node := range nodes{
		for _, connection := range node.connections{
			connection.back.weight += node.weight * connection.value;
		}
	}
}
func activateNodes(nodes [] Neuron){
	for _ , node := range nodes{
		node.nonActivatedWeight = node.weight;
		node.weight = sigmoidFunc(node.weight);
	}
}

func backwardProp(inNodes [] Neuron,hidNodes [] Neuron,outNodes [] Neuron, deltaOutSum float64){
	// since this for a xor function..
	// input is a const
	// var expectedResult := strconv.ParseFloat(input[0], 64)^strconv.ParseFloat(input[1], 64);
	// fmt.Println("expectedResult = %f", expectedResult);
	// var calculated := outNodes[0]; 
	// var sumMarginOfError := expectedResult - calculated;

	// var deltaOutputSum := derivativeSigmoid(calculated) * sumMarginOfError;
	var deltaWeight := [len(hidNodes)] int;				// TODO!!!!! MAGIC
	var outNode := outNodes[0];



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

func runSimulation(in [] Neuron, hid [] Neuron, out [] Neuron){
	// run through process
	forwardSim(in,hid,out);
	backwardSim(in,hid,out);
	// forwardProp(in);
	// activateNodes(hid);
	// forwardProp(hid);
	// activateNodes(out);
	fmt.Println("FINAL == %f", out[0].weight);


}
func forwardSim(in [] Neuron, hid [] Neuron, out [] Neuron){
	forwardProp(in);
	activateNodes(hid);
	forwardProp(hid);
	activateNodes(out);
}
// this should also probably take in or produce the delta output sum value
func backwardSim(in [] Neuron, hid [] Neuron, out [] Neuron){

	var expectedResult := strconv.ParseFloat(input[0], 64)^strconv.ParseFloat(input[1], 64);
	fmt.Println("expectedResult = %f", expectedResult);
	var calculated := outNodes[0]; 
	var sumMarginOfError := expectedResult - calculated;

	var deltaOutputSum := derivativeSigmoid(calculated) * sumMarginOfError
	backwardProp(in,hid,out, deltaOutputSum);
}
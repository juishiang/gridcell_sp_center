package main

import (
	"image/color"
	"math"
	"math/rand"
	"strconv"

	"fmt"
	"os"
	"time"

	gridcell "github.com/juishiang/GridCell"
	"gonum.org/v1/gonum/stat"

	//"golang.org/x/exp/errors/fmt"
	/// gonum plot:
	/*
		"gonum.org/v1/plot"
		"gonum.org/v1/plot/plotter"
		"gonum.org/v1/plot/plotutil"
		"gonum.org/v1/plot/vg"
	*/
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	///// mode parameter
	//totaltimeS := []int{10000000, 100000000, 1000000000}
	probresS := []float64{0.001} //[]float64{0.001, 0.0001}
	//totaltime * probres >= 10^6 line is clear
	//var probres float64 = 0.0001
	plus := []bool{true, false}
	var P bool
	var cosmode bool = true
	///// functional parameter
	var mode string
	var diff float64 = 0.0
	var filename string
	///// mode parameter
	var simDat gridcell.Gridcellcos /// for cos mode cell
	var simDat2 gridcell.Grid_cell  /// for normal mode
	////2 central percent
	spacingthresholdS := []float32{0.3, 0.4, 0.5, 0.6, 0.7}
	var spacingthreshold float32 = 0.4
	////spacing setting
	var spacingsize float64 = 100.0
	var centralspacing1 float64 = 50.0
	var centralspacing2 float64 = 40.0
	////orientation setting
	var centralori1 float64 = -15.0
	var centralori2 float64 = 15.0
	var spacingU float64 = 50.0
	//spacingVS := []float64{3.0, 5.0, 7.0, 10.0}
	spacingVS := []float64{3.0}
	var spacingV float64
	////theta setting
	var thetaU float64 = 0.0
	//thetaVS := []float64{3.0, 5.0, 7.0, 10.0}
	thetaVS := []float64{3.0}
	var thetaV float64
	var positionC float64 = 0.5 * spacingsize //uniform distribution from -0.5 spacing to 0.5 spacing [-0.5*spacing 0.5*spacing]
	var statRes float64 = 0.5
	s := time.Now()
	var avgentropym1, avgentropym2 float64 = 0.0, 0.0
	for k := range probresS {
		for _, spacingthreshold = range spacingthresholdS {
			probres := probresS[k]
			totaltime := int(1.e+6 / probres)
			karr := make([]int, int(1.0/probres)+1)
			karrf := make([]float64, int(1.0/probres)+1)
			for _, spacingV = range spacingVS {
				for _, thetaV = range thetaVS {
					if cosmode {
						avgentropym1 = 0.0
						avgentropym2 = 0.0
						for avgtime := 0; avgtime < 10; avgtime++ {
							for _, P = range plus {
								for i := range karr {
									karr[i] = 0
									karrf[i] = 0.0
								}
								rand.Seed(time.Now().UnixNano())
								spacingDi := make([]int, int(float64(spacingsize)/statRes)+1)
								oriDi := make([]int, int(120.0/statRes)+1)
								////spacing, orientation distribution
								pointSO := plotter.XYs{}
								pointS := plotter.XYs{}
								pointO := plotter.XYs{}
								////end
								for i := 0; i < totaltime; i++ {
									//// setting grid cell (cos mode) parameter
									if rand.Float32() > spacingthreshold {
										spacingU = centralspacing1
										thetaU = centralori1
									} else {
										spacingU = centralspacing2
										thetaU = centralori2
									}
									spacing := spacingU + rand.NormFloat64()*spacingV
									theta := thetaU + rand.NormFloat64()*thetaV
									posix := positionC + spacing*(rand.Float64()-0.5)
									posiy := positionC + spacing*(rand.Float64()-0.5)
									simDat.Init(spacing, theta, 0.0, 0.0, 0.0)
									karr[int(simDat.Activation(posix, posiy, P)/probres)]++
									oriDi[int((theta+60.0)/statRes)]++
									spacingDi[int((spacing)/statRes)]++
									////Spacing and orientation distribution record
									pointSO = append(pointSO, plotter.XY{
										X: theta,
										Y: spacing,
									})
								}
								//// gonum plot name
								if P {
									mode = "plus"
								} else {
									mode = "product"
								}
								filename = "pdf_1e-3" + mode + "TV_" + strconv.Itoa(int(thetaV)) + "_SV_" + strconv.Itoa(int(spacingV)) + ".png"

								str, _ := os.Getwd()
								folder := "./sp_thre=" + strconv.Itoa(int(spacingthreshold*10)) + "/"
								err := os.MkdirAll(folder, os.ModePerm)
								filename = folder + filename
								p, err := plot.New()
								if err != nil {
									panic(err)
								}
								p.Y.Min, p.X.Min, p.Y.Max, p.X.Max = 20.0, -30.0, 60.0, 30.0
								p.Title.Text = "spacing vs orientation"
								p.X.Label.Text = "orientation"
								p.Y.Label.Text = "spacing"
								// Draw a grid behind the data
								p.Add(plotter.NewGrid())

								// Make a scatter plotter and set its style.
								ss, err := plotter.NewScatter(pointSO)
								if err != nil {
									panic(err)
								}
								ss.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
								p.Add(ss)

								// Save the plot to a PNG file.
								if err := p.Save(4*vg.Inch, 4*vg.Inch, str+"SVSOscatter.png"); err != nil {
									panic(err)
								}
								for idx := range spacingDi {
									//spacingDf[idx] = float64(spacingDi[idx])/(float64(neunum)*statRes)
									pointS = append(pointS, plotter.XY{
										X: float64(idx) * statRes,
										Y: float64(spacingDi[idx]),
									})
								}
								gonumplt, err := plot.New()
								if err != nil {
									panic(err)
								}
								gonumplt.Title.Text = "spacing"
								gonumplt.X.Label.Text = "spacing"
								gonumplt.Y.Label.Text = "count"
								//gonumplt.Y.Min, gonumplt.X.Min, gonumplt.Y.Max, gonumplt.X.Max = 0.0, 0.0, 1.0, 1.0
								if err := plotutil.AddLines(gonumplt,
									"spacing distribution", pointS,
								); err != nil {
									panic(err)
								}
								if err := gonumplt.Save(5*vg.Inch, 5*vg.Inch, str+"spacingdistri.png"); err != nil {
									panic(err)
								}
								for idx := range oriDi {
									//oriDf[idx] = float64(oriDi[idx])/(float64(neunum)*statRes)
									pointO = append(pointO, plotter.XY{
										X: float64(idx) * statRes,
										Y: float64(oriDi[idx]),
									})
								}
								gonumplt2, err := plot.New()
								if err != nil {
									panic(err)
								}
								gonumplt2.Title.Text = "orientation"
								gonumplt2.X.Label.Text = "degree"
								gonumplt2.Y.Label.Text = "count"
								//gonumplt.Y.Min, gonumplt.X.Min, gonumplt.Y.Max, gonumplt.X.Max = 0.0, 0.0, 1.0, 1.0
								if err := plotutil.AddLines(gonumplt2,
									"spacing distribution", pointO,
								); err != nil {
									panic(err)
								}
								if err := gonumplt2.Save(5*vg.Inch, 5*vg.Inch, str+"orientation_distri.png"); err != nil {
									panic(err)
								}
								if cosmode {
									if P {
										diff = Pdfpltentr(karr, karrf, probres, mode, filename, totaltime)
										avgentropym1 += diff
									} else {
										enr := Pdfpltentr(karr, karrf, probres, mode, filename, totaltime)
										diff -= enr
										avgentropym2 += enr
										fmt.Printf("plus mode value is larger than product mode with %f bits(log2)\n", diff)
									}
								}
								fmt.Println("exc time:", time.Since(s))
								s = time.Now()
							}
						}
						avgentropym1 /= 10.0
						avgentropym2 /= 10.0
						fmt.Printf("avg entropy (10times) in TV %d SV %d \n plus: %f \n product: %f \n ", int(thetaV), int(spacingV), avgentropym1, avgentropym2)
						avgentropym1 = 0.0
						avgentropym2 = 0.0
					} else {
						for i := range karr {
							karr[i] = 0
							karrf[i] = 0.0
						}
						rand.Seed(time.Now().UnixNano())
						for i := 0; i < totaltime; i++ {
							//// setting grid cell (cos mode) parameter
							spacing := spacingU + rand.NormFloat64()*spacingV
							if spacing <= 0.1 {
								i--
								continue
							}
							theta := thetaU + rand.NormFloat64()*thetaV
							posix := positionC + spacing*(rand.Float64()-0.5)
							posiy := positionC + spacing*(rand.Float64()-0.5)
							//fmt.Println(spacing, int(spacingsize/spacing)+1, posix, posiy, theta)
							simDat2.Init(int(spacingsize/spacing)+1, 0.0, 0.0, spacingsize, 1.0, 0.0, 0.2)
							simDat2.AcK(spacing, 0.0, 0.0, 0.3)
							//fmt.Println(simDat2.Fireact(posix, posiy, 0.0))
							karr[int(simDat2.Fireact(posix, posiy, theta)/probres)]++
						}
						filename = "pdf_1e-3" + "normal_distr" + ".png"
						mode = "normal distr"
						Pdfpltentr(karr, karrf, probres, mode, filename, totaltime)
						fmt.Println("exc time:", time.Since(s))
						s = time.Now()
					}
				}
			}
		}
	}
}

//Pdfpltentr is used to calculate the pdf entropy and plt the distribution
//gonum plot
//fmt.Println("start save fig")
func Pdfpltentr(karr []int, karrf []float64, probres float64, mode, filename string, totaltime int) (entr float64) {
	points := plotter.XYs{}
	for i := range karr {
		points = append(points, plotter.XY{
			X: float64(i) * probres,
			Y: float64(karr[i]) / (float64(totaltime) * probres),
		})
		karrf[i] = float64(karr[i]) / (float64(totaltime))
	}
	plt, err := plot.New()
	if err != nil {
		panic(err)
	}
	plt.Y.Min, plt.X.Min, plt.Y.Max, plt.X.Max = 0.0, 0.0, 1.0, 1.0

	if err := plotutil.AddLines(plt,
		"line1", points,
	); err != nil {
		panic(err)
	}
	if err := plt.Save(5*vg.Inch, 5*vg.Inch, filename); err != nil {
		panic(err)
	}
	//fmt.Println("pic save, start cal entropy")
	entr = (stat.Entropy(karrf) / math.Ln2) + math.Log2(probres)
	fmt.Printf("entropy of %s in resolution %f is : %f\n", filename, probres, entr)
	////test prob sum = 1
	/*
		var sumT float64 = 0.0
		for l := range karrf {
			sumT += karrf[l]
		}
		fmt.Printf("chech sum: sum of prob = %f\n", sumT)
	*/
	return
}

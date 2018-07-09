# External Sorting on Goroutine

This is a small library that implements External Merge Sort in Golang

    "External sorting is a term for a class of sorting algorithms that can handle massive amounts of data. External sorting is required when the data being sorted do not fit into the main memory of a computing device (usually RAM) and instead they must reside in the slower external memory (usually a hard drive). External sorting typically uses a hybrid sort-merge strategy. In the sorting phase, chunks of data small enough to fit in main memory are read, sorted, and written out to a temporary file. In the merge phase, the sorted subfiles are combined into a single larger file. External sorting is a form of distribution sort, with the added characteristic that the individual subsets are separately sorted, rather than just being used internally as intermediate stages." – Wikipedia

[TOCM]

[TOC]

##Introduction

The library is pretty flexible as it takes number of elements and empty file's name  as input to store the initial big data and returns sorted result in a new file.

Sort order is default by sort.Interface.

Persistence is handled by an implementation of the package pipeline, of which there are currently four implementations:

[GenerateRandomFile](https://github.com/vergilken/external_sort_with_goroutine/blob/master/pipeline/nodes.go#L129) - Generating elements and sort them in a empty file
[CreatePipeline](https://github.com/vergilken/external_sort_with_goroutine/blob/master/pipeline/nodes.go#L142) - Implementing channels to control goroutine function  to manipulate elements sorting in memory and out of memory
[WriteTofile](https://github.com/vergilken/external_sort_with_goroutine/blob/master/pipeline/nodes.go#L163) - Wrting results to the new file with different goroutine pipeline 

##Test Results

######100000000 random integers(64 bit) in 4 chunks testing:

Filename: large, Size: 800000000
Read Done:  3.771786175s
Read Done:  4.364625052s
Read Done:  5.290082783s
Read Done:  5.296543703s
In Memory Sort Done:  10.1021381s
In Memory Sort Done:  10.557499954s
In Memory Sort Done:  11.40966397s
In Memory Sort Done:  11.431722567s
Merge Done:  30.469568818s
Merge Done:  30.469595596s
Merge Done:  30.469794059s
Filename: Sorted_large, Size: 800000000

// MIT License

// Copyright (c) 2024 ISSuh

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package message

import (
	"github.com/ISSuh/sos/internal/domain/model/entity"
	"github.com/ISSuh/sos/pkg/empty"
	"github.com/ISSuh/sos/pkg/validation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromObjectID(objectID entity.ObjectID) *ObjectID {
	return &ObjectID{
		Id: objectID.ToInt64(),
	}
}

func ToObjectID(objectID *ObjectID) entity.ObjectID {
	if validation.IsNil(objectID) {
		return empty.Struct[entity.ObjectID]()
	}
	return entity.NewObjectIDFrom(objectID.Id)
}

func FromBlockID(blockID entity.BlockID) *BlockID {
	return &BlockID{
		Id: blockID.ToInt64(),
	}
}

func ToBlockID(blockID *BlockID) entity.BlockID {
	if validation.IsNil(blockID) {
		return empty.Struct[entity.BlockID]()
	}
	return entity.NewBlockIDFrom(blockID.Id)
}

func FromObjectMetadata(objectMetadata entity.ObjectMetadata) *ObjectMetadata {
	return &ObjectMetadata{
		Id:        FromObjectID(objectMetadata.ID()),
		Group:     objectMetadata.Group(),
		Partition: objectMetadata.Partition(),
		Path:      objectMetadata.Path(),
		Name:      objectMetadata.Name(),
		Size:      int32(objectMetadata.Size()),
	}
}

func ToObjectMetadata(objectMetadata *ObjectMetadata) entity.ObjectMetadata {
	if validation.IsNil(objectMetadata) {
		return empty.Struct[entity.ObjectMetadata]()
	}

	builder := entity.NewObjectMetadataBuilder()
	builder.ID(ToObjectID(objectMetadata.Id)).
		Group(objectMetadata.Group).
		Partition(objectMetadata.Partition).
		Path(objectMetadata.Path).
		Name(objectMetadata.Name).
		Size(int(objectMetadata.Size))

	return builder.Build()
}

func FromBlockHeader(blockHeader entity.BlockHeader) *BlockHeader {
	return &BlockHeader{
		ObjectID:  FromObjectID(blockHeader.ObjectID()),
		BlockID:   FromBlockID(blockHeader.BlockID()),
		Index:     int32(blockHeader.Index()),
		Size:      int32(blockHeader.Size()),
		Checksum:  blockHeader.Checksum(),
		Timestamp: timestamppb.New(blockHeader.Timestamp()),
	}
}

func ToBlockHeader(blockHeader *BlockHeader) entity.BlockHeader {
	if validation.IsNil(blockHeader) {
		return empty.Struct[entity.BlockHeader]()
	}

	builder := entity.NewBlockHeaderBuilder()
	builder.
		ObjectID(ToObjectID(blockHeader.ObjectID)).
		BlockID(ToBlockID(blockHeader.BlockID)).
		Index(int(blockHeader.Index)).
		Size(int(blockHeader.Size)).
		Timestamp(blockHeader.Timestamp.AsTime())

	return builder.Build()
}

func FromBlock(block entity.Block) *Block {
	return &Block{
		Header: FromBlockHeader(block.Header()),
		Data:   block.Buffer(),
	}
}

func ToBlock(block *Block) entity.Block {
	if validation.IsNil(block) {
		return empty.Struct[entity.Block]()
	}

	builder := entity.NewBlockBuilder()
	builder.Header(ToBlockHeader(block.Header)).
		Buffer(block.Data)

	return builder.Build()
}

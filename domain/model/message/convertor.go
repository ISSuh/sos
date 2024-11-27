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
	"github.com/ISSuh/sos/domain/model/dto"
	"github.com/ISSuh/sos/domain/model/entity"
	"github.com/ISSuh/sos/internal/empty"
	"github.com/ISSuh/sos/internal/validation"
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

func FromVersionDTO(version *dto.Version) *Version {
	blockHeaders := make([]*BlockHeader, 0, len(version.BlockHeaders))
	for _, header := range version.BlockHeaders {
		blockHeaders = append(blockHeaders, FromBlockHeaderDTO(&header))
	}

	return &Version{
		Number:       int32(version.Number),
		Size:         int32(version.Size),
		BlockHeaders: blockHeaders,
		CreatedAt:    timestamppb.New(version.CreatedAt),
		ModifiedAt:   timestamppb.New(version.ModifiedAt),
	}
}

func ToVersionDTO(version *Version) *dto.Version {
	if validation.IsNil(version) {
		return nil
	}

	blockHeaders := make([]dto.BlockHeader, 0, len(version.BlockHeaders))
	for _, header := range version.BlockHeaders {
		blockHeaders = append(blockHeaders, ToBlockHeaderDTO(header))
	}

	return &dto.Version{
		Number:       int(version.Number),
		Size:         int(version.Size),
		BlockHeaders: blockHeaders,
		CreatedAt:    version.CreatedAt.AsTime(),
		ModifiedAt:   version.ModifiedAt.AsTime(),
	}
}

func FromObjectMetadata(objectMetadata entity.ObjectMetadata) *ObjectMetadata {
	return &ObjectMetadata{
		Id:         FromObjectID(objectMetadata.ID()),
		Group:      objectMetadata.Group(),
		Partition:  objectMetadata.Partition(),
		Path:       objectMetadata.Path(),
		Name:       objectMetadata.Name(),
		CreatedAt:  timestamppb.New(objectMetadata.CreatedAt),
		ModifiedAt: timestamppb.New(objectMetadata.ModifiedAt),
	}
}

func FromObjectMetadataDTO(objectMetadata *dto.Metadata) *ObjectMetadata {
	versions := make([]*Version, 0, len(objectMetadata.Versions))
	for _, version := range objectMetadata.Versions {
		versions = append(versions, FromVersionDTO(&version))
	}

	return &ObjectMetadata{
		Id:         FromObjectID(objectMetadata.ID),
		Group:      objectMetadata.Group,
		Partition:  objectMetadata.Partition,
		Path:       objectMetadata.Path,
		Name:       objectMetadata.Name,
		Versions:   versions,
		CreatedAt:  timestamppb.New(objectMetadata.CreatedAt),
		ModifiedAt: timestamppb.New(objectMetadata.ModifiedAt),
	}
}

func ToObjectMetadataDTO(objectMetadata *ObjectMetadata) *dto.Metadata {
	if validation.IsNil(objectMetadata) {
		return nil
	}

	versions := make([]dto.Version, 0, len(objectMetadata.Versions))
	for _, version := range objectMetadata.Versions {
		versions = append(versions, *ToVersionDTO(version))
	}

	return &dto.Metadata{
		ID:         ToObjectID(objectMetadata.Id),
		Group:      objectMetadata.Group,
		Partition:  objectMetadata.Partition,
		Name:       objectMetadata.Name,
		Versions:   versions,
		Path:       objectMetadata.Path,
		CreatedAt:  objectMetadata.CreatedAt.AsTime(),
		ModifiedAt: objectMetadata.ModifiedAt.AsTime(),
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
		Name(objectMetadata.Name)

	return builder.Build()
}

func FromObjectMetadataListDTO(list dto.MetadataList) *ObjectMetadataList {
	metadataList := make([]*ObjectMetadata, 0, len(list))
	for _, metadata := range list {
		metadataList = append(metadataList, FromObjectMetadataDTO(&metadata))
	}

	return &ObjectMetadataList{
		Metadata: metadataList,
	}
}

func ToObjectMetadataListDTO(list *ObjectMetadataList) dto.MetadataList {
	metadataList := make(dto.MetadataList, 0, len(list.Metadata))
	for _, metadata := range list.Metadata {
		metadataList = append(metadataList, *ToObjectMetadataDTO(metadata))
	}

	return metadataList
}

func ToItemsDTO(list *ObjectMetadataList) dto.Items {
	metadataList := make(dto.MetadataList, 0, len(list.Metadata))
	for _, metadata := range list.Metadata {
		metadataList = append(metadataList, *ToObjectMetadataDTO(metadata))
	}

	return dto.NewItemsFromMetadataList(metadataList)
}

func FromBlockHeader(blockHeader *entity.BlockHeader) *BlockHeader {
	return &BlockHeader{
		ObjectID:  FromObjectID(blockHeader.ObjectID()),
		BlockID:   FromBlockID(blockHeader.BlockID()),
		Index:     int32(blockHeader.Index()),
		Size:      int32(blockHeader.Size()),
		Checksum:  blockHeader.Checksum(),
		Timestamp: timestamppb.New(blockHeader.Timestamp()),
	}
}

func FromBlockHeaderDTO(blockHeader *dto.BlockHeader) *BlockHeader {
	return &BlockHeader{
		ObjectID:  FromObjectID(blockHeader.ObjectID),
		BlockID:   FromBlockID(blockHeader.BlockID),
		Index:     int32(blockHeader.Index),
		Size:      int32(blockHeader.Size),
		Checksum:  blockHeader.Checksum,
		Timestamp: timestamppb.New(blockHeader.Timestamp),
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
		Checksum(blockHeader.Checksum).
		Timestamp(blockHeader.Timestamp.AsTime())

	return builder.Build()
}

func ToBlockHeaderDTO(blockHeader *BlockHeader) dto.BlockHeader {
	if validation.IsNil(blockHeader) {
		return empty.Struct[dto.BlockHeader]()
	}
	return dto.BlockHeader{
		ObjectID:  ToObjectID(blockHeader.ObjectID),
		BlockID:   ToBlockID(blockHeader.BlockID),
		Index:     int(blockHeader.Index),
		Size:      int(blockHeader.Size),
		Checksum:  blockHeader.Checksum,
		Timestamp: blockHeader.Timestamp.AsTime(),
	}
}

func FromBlock(block *entity.Block) *Block {
	header := block.Header()
	return &Block{
		Header: FromBlockHeader(&header),
		Data:   block.Buffer(),
	}
}

func FromBlockDTO(block *dto.Block) *Block {
	header := FromBlockHeaderDTO(&block.Header)
	return &Block{
		Header: header,
		Data:   block.Data,
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

func FromObjectDTO(object *dto.Object) *Object {
	blockHeaders := make([]*BlockHeader, 0, len(object.BlockHeaders))
	for _, header := range object.BlockHeaders {
		blockHeaders = append(blockHeaders, FromBlockHeaderDTO(&header))
	}

	return &Object{
		Id:           FromObjectID(object.ID),
		Group:        object.Group,
		Partition:    object.Partition,
		Name:         object.Name,
		Path:         object.Path,
		Size:         int32(object.Size),
		VersionNum:   int32(object.VersionNum),
		BlockHeaders: blockHeaders,
	}
}

func ToObjectDTO(object *Object) *dto.Object {
	if validation.IsNil(object) {
		return nil
	}

	blockHeaders := make([]dto.BlockHeader, 0, len(object.BlockHeaders))
	for _, header := range object.BlockHeaders {
		blockHeaders = append(blockHeaders, ToBlockHeaderDTO(header))
	}

	return &dto.Object{
		ID:           ToObjectID(object.Id),
		Group:        object.Group,
		Partition:    object.Partition,
		Name:         object.Name,
		Path:         object.Path,
		Size:         int(object.Size),
		VersionNum:   int(object.VersionNum),
		BlockHeaders: blockHeaders,
	}
}

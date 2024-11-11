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

package object

import (
	"context"

	"github.com/ISSuh/sos/internal/domain/model/dto"
	"github.com/ISSuh/sos/internal/domain/model/message"
	"github.com/ISSuh/sos/internal/infrastructure/transport/rpc"
)

type Deleter struct {
	objectRequestor  rpc.MetadataRegistryRequestor
	storageRequestor rpc.BlockStorageRequestor
}

func NewDeleter(objectRequestor rpc.MetadataRegistryRequestor, storageRequestor rpc.BlockStorageRequestor) Deleter {
	return Deleter{
		objectRequestor:  objectRequestor,
		storageRequestor: storageRequestor,
	}
}

func (o *Deleter) Delete(c context.Context, metadata dto.Metadata) error {
	for _, version := range metadata.Versions {
		if err := o.deleteBlocks(c, version); err != nil {
			return err
		}
	}

	metadata.Versions = nil
	if err := o.deleteObjectMetadata(c, metadata); err != nil {
		return err
	}

	return nil
}

func (o *Deleter) DeleteVersion(c context.Context, metadata dto.Metadata, deleteVersionNum int) error {
	version, err := metadata.Versions.Version(deleteVersionNum)
	if err != nil {
		return err
	}

	if err := o.deleteBlocks(c, version); err != nil {
		return err
	}

	metadata.Versions = dto.Versions{
		version,
	}

	if err := o.deleteObjectMetadata(c, metadata); err != nil {
		return err
	}

	return nil
}

func (o *Deleter) deleteObjectMetadata(c context.Context, metadata dto.Metadata) error {
	msg := metadata.ToMessage()
	if _, err := o.objectRequestor.Delete(c, msg); err != nil {
		return err
	}

	return nil
}

func (o *Deleter) deleteBlocks(c context.Context, version dto.Version) error {
	blockHeaders := version.BlockHeaders
	for _, blockHeader := range blockHeaders {
		msg := &message.BlockHeader{
			ObjectID: &message.ObjectID{
				Id: blockHeader.ObjectID.ToInt64(),
			},
			BlockID: &message.BlockID{
				Id: blockHeader.BlockID.ToInt64(),
			},
			Index: int32(blockHeader.Index),
		}

		if _, err := o.storageRequestor.Delete(c, msg); err != nil {
			return err
		}
	}

	return nil
}

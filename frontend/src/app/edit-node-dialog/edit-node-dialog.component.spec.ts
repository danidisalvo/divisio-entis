import {ComponentFixture, TestBed} from '@angular/core/testing';
import {Color} from '@angular-material-components/color-picker';
import {FormBuilder} from '@angular/forms'
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";

import {EditNodeDialogComponent} from './edit-node-dialog.component';

describe('EditNodeDialogComponent', () => {
  let component: EditNodeDialogComponent;
  let fixture: ComponentFixture<EditNodeDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      providers: [
        {
          provide: FormBuilder
        },
        {
          provide: MatDialogRef,
          useValue: {}
        },
        {
          provide: MAT_DIALOG_DATA,
          useValue: {
            node: {
              name: 'name',
              color: new Color(0, 0, 0),
            },
            child: {
              name: '',
              color: new Color(0, 0, 0),
            },
            addChild: false
          }
        }
      ],
      declarations: [EditNodeDialogComponent]
    })
      .compileComponents();

    fixture = TestBed.createComponent(EditNodeDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
